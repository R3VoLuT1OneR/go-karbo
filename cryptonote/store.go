package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/r3volut1oner/go-karbo/config"
)

var (
	ErrBlockExists = errors.New("block already exists")

	ErrBlockNotFound = errors.New("block not found")
)

type Store interface {
	// Init database files with genesis block if empty.
	Init(n *config.Network) error

	// GetBlockHashByHeight provides block by hash
	GetBlockHashByHeight(uint32) (*Hash, error)

	// GetBlockHeightByHash returns height for specific block hash
	GetBlockIndexByHash(*Hash) (uint32, error)

	// GetBlockByHeight returns block by height
	GetBlockByHeight(uint32) (*Block, error)

	// AppendBlock to database persistence layer.
	AppendBlock(*Block) error

	// HasBlock verifies that block is saved in DB
	HasBlock(*Hash) (bool, error)

	// GetHeight of current saved blockchain
	TopIndex() (uint32, error)

	// Empty checks if database is new and empty
	Empty() (bool, error)

	// Close
	Close() error
}

type badgerDB struct {
	badger *badger.DB
}

type storeTxn struct {
	*badger.Txn
}

func itob(i uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], i)
	return buf[:]
}

func btoi(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func keyBlockByHash(h *Hash) []byte {
	return []byte("block-" + h.String())
}

func keyBlockHeightByHash(h *Hash) []byte {
	return []byte("block-height-" + h.String())
}

func keyBlockHashByHeight(h uint32) []byte {
	return append([]byte("block-index-"), itob(uint64(h))...)
}

func keyHeight() []byte {
	return []byte("block-height")
}

func NewBadgerDB() (Store, error) {
	db, err := badger.Open(
		badger.DefaultOptions("./.badger"),
	)

	if err != nil {
		return nil, err
	}

	return &badgerDB{
		badger: db,
	}, nil
}

func (db *badgerDB) Init(n *config.Network) error {
	block, err := GenerateGenesisBlock(n)
	if err != nil {
		return err
	}

	genesisHash, err := block.Hash()
	if err != nil {
		return err
	}

	height, err := db.GetBlockIndexByHash(genesisHash)

	if err == ErrBlockNotFound {
		if err := db.AppendBlock(block); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if height != 0 {
		return errors.New("genesis block height is not match")
	}

	return nil
}

func (db *badgerDB) AppendBlock(b *Block) error {
	hash, err := b.Hash()
	if err != nil {
		return err
	}

	payload, err := b.Serialize()
	if err != nil {
		return err
	}

	return db.badger.Update(func(txn *badger.Txn) error {
		keyHeight := keyHeight()
		keyBlock := keyBlockByHash(hash)
		keyHashByHeight := keyBlockHeightByHash(hash)

		if _, err := txn.Get(keyBlock); err == nil {
			return ErrBlockExists
		}

		blockHeight := uint64(0)
		heightItem, err := txn.Get(keyHeight)

		if err != nil && err != badger.ErrKeyNotFound {
			return err
		} else if err == nil {
			heightPayload, err := heightItem.ValueCopy(nil)
			if err != nil {
				return err
			}

			blockHeight = btoi(heightPayload) + 1
		}

		keyBlockHeight := keyBlockHashByHeight(uint32(blockHeight))

		if err := txn.Set(keyBlock, payload); err != nil {
			return err
		}

		if err := txn.Set(keyHeight, itob(blockHeight)); err != nil {
			return err
		}

		if err := txn.Set(keyHashByHeight, itob(blockHeight)); err != nil {
			return err
		}

		if err := txn.Set(keyBlockHeight, hash[:]); err != nil {
			return err
		}

		return nil
	})
}

func (db *badgerDB) HasBlock(hash *Hash) (bool, error) {
	hasBlock := false

	err := db.badger.View(func(txn *badger.Txn) error {
		_, err := txn.Get(keyBlockByHash(hash))

		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		hasBlock = err != badger.ErrKeyNotFound

		return nil
	})

	if err != nil {
		return false, err
	}

	return hasBlock, nil
}

func (db *badgerDB) GetBlockHashByHeight(height uint32) (*Hash, error) {
	var h *Hash

	err := db.badger.View(func(txn *badger.Txn) error {
		stxn := &storeTxn{txn}

		foundHash, err := stxn.getHashByKey(keyBlockHashByHeight(height))
		if err != nil {
			return err
		}

		h = foundHash

		return nil
	})

	if err != nil {
		return nil, err
	}

	return h, err
}

func (db *badgerDB) GetBlockIndexByHash(h *Hash) (uint32, error) {
	var height uint32

	err := db.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(keyBlockHeightByHash(h))

		if err == badger.ErrKeyNotFound {
			return ErrBlockNotFound
		} else if err != nil {
			return err
		}

		hb, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		height = uint32(btoi(hb))

		return nil
	})

	if err != nil {
		return 0, err
	}

	return height, nil
}

func (db *badgerDB) GetBlockByHeight(height uint32) (*Block, error) {
	var b *Block

	err := db.badger.View(func(txn *badger.Txn) error {
		stxn := &storeTxn{txn}

		foundHash, err := stxn.getHashByKey(keyBlockHashByHeight(height))
		if err != nil {
			return err
		}

		block, err := stxn.getBlockByKey(keyBlockByHash(foundHash))
		if err != nil {
			return err
		}

		b = block

		return nil
	})

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (db *badgerDB) Empty() (bool, error) {
	height, err := db.TopIndex()
	if err != nil {
		return false, err
	}

	return height == 0, nil
}

func (db *badgerDB) Close() error {
	if err := db.badger.Close(); err != nil {
		return err
	}

	return nil
}

func (db *badgerDB) TopIndex() (uint32, error) {
	height := uint64(0)

	err := db.badger.View(func(txn *badger.Txn) error {
		keyHeight := keyHeight()

		heightItem, err := txn.Get(keyHeight)
		if err != nil {
			return err
		}

		heightPayload, err := heightItem.ValueCopy(nil)
		if err != nil {
			return err
		}

		height = btoi(heightPayload)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return uint32(height), nil
}

func (txn *storeTxn) getHashByKey(b []byte) (*Hash, error) {
	var h Hash

	item, err := txn.Get(b)

	if err == badger.ErrKeyNotFound {
		return nil, ErrBlockNotFound
	} else if err != nil {
		return nil, err
	}

	hb, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	if err := h.Read(bytes.NewReader(hb)); err != nil {
		return nil, err
	}

	return &h, nil
}

func (txn *storeTxn) getBlockByKey(b []byte) (*Block, error) {
	var block Block

	item, err := txn.Get(b)

	if err == badger.ErrKeyNotFound {
		return nil, ErrBlockNotFound
	} else if err != nil {
		return nil, err
	}

	bb, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	if err := block.Deserialize(bytes.NewReader(bb)); err != nil {
		return nil, err
	}

	return &block, nil
}