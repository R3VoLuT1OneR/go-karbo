package cryptonote

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type BlockChain struct {
	// Network is current network configurations, must stay immutable
	Network *config.Network

	// logger for block chain events
	logger *log.Logger

	// genesisBlock used for caching genesis block
	genesisBlock *Block

	blocksIndex map[Hash]*Block

	sync.RWMutex
}

// NewBlockChain generates basic blockchain object
func NewBlockChain(network *config.Network, logger *log.Logger) *BlockChain {
	bc := &BlockChain{
		Network: network,
		logger:  logger,
	}

	return bc
}

// AddBlock used for adding new blocks to the blockchain.
//
// It returns nil if block added successfully and ErrAddBlock* in case of error
func (bc *BlockChain) AddBlock(b *Block, rawTransactions [][]byte) error {
	bc.Lock()
	defer bc.Unlock()

	hash := b.Hash()
	height := b.Height()

	blogger := bc.logger.WithFields(log.Fields{
		"block_hash":   hash,
		"block_height": height,
	})

	if bc.haveBlock(hash) {
		err := ErrAddBlockAlreadyExists
		blogger.Error(err)
		return err
	}

	if !bc.haveBlock(&b.PreviousBlockHash) {
		err := ErrAddBlockRejectedAsOrphaned
		blogger.Error(err)
		return err
	}

	coinbaseTransactionSize := b.CoinbaseTransaction.Size()
	if coinbaseTransactionSize > bc.Network.MaxTxSize {
		err := ErrAddBlockTransactionCoinbaseMaxSize
		blogger.Error(err)
		return err
	}

	if len(b.TransactionsHashes) != len(rawTransactions) {
		err := ErrAddBlockTransactionCountNotMatch
		blogger.Error(err)
		return err
	}

	var transactions []Transaction
	var transactionsSize uint64
	if transactions, transactionsSize, err := bc.deserializeTransactions(blogger, rawTransactions); err != nil {
		return err
	}

	prevBlockHeight := bc.blockHeight(&b.PreviousBlockHash)

	blockSize := coinbaseTransactionSize + transactionsSize
	if blockSize > bc.Network.MaxBlockSize(prevBlockHeight) {
		err := ErrBlockValidationCumulativeSizeTooBig
		blogger.Error(err)
		return err
	}

	var minerReward uint64
	if minerReward, err := bc.validateBlock(b); err != nil {
		blogger.Error(err)
		return err
	}

	return nil
}

// validateBlock validates block
//
// Returns miner reward and an error if there is an error in block validation
func (bc *BlockChain) validateBlock(blogger *log.Entry, b *Block) (uint64, error) {
	prevBlockHeight := bc.blockHeight(&b.PreviousBlockHash)
	minerReward := 0

	if bc.Network.GetBlockMajorVersion(b.Height()) != b.MajorVersion {
		return 0, ErrBlockValidationWrongVersion
	}

	if b.MajorVersion == config.BlockMajorVersion2 && b.Parent.MajorVersion > config.BlockMajorVersion1 {
		err := ErrBlockValidationParentBlockWrongVersion
		blogger.WithField("block_parent_major_version", b.Parent.MajorVersion).Error(err)
		return 0, err
	}

	if b.MajorVersion == config.BlockMajorVersion2 || b.MajorVersion == config.BlockMajorVersion3 {
		if len(b.Parent.serialize(false)) > 2048 {
			err := ErrBlockValidationParentBlockSizeTooBig
			blogger.Error(err)
			return 0, err
		}
	}

	if b.Timestamp > uint64(time.Now().Unix())+bc.Network.BlockFutureTimeLimit(b.MajorVersion) {
		err := ErrBlockValidationTimestampTooFarInFuture
		blogger.Error(err)
		return 0, err
	}

}

// blockHeight returns index on the current block
func (bc *BlockChain) blockHeight(h *Hash) uint64 {
	return 0
}

// TODO: Properly implement this method
// haveBlock return whether the block hash contains in the blockchain
//
// This function is NOT safe for concurrent access
func (bc *BlockChain) haveBlock(h *Hash) bool {
	_, ok := bc.blocksIndex[*h]
	return ok
}

// HaveBlock return whether the block hash contains in the blockchain
//
// This function is safe for concurrent access.
func (bc *BlockChain) HaveBlock(h *Hash) bool {
	bc.RLock()
	hasBlock := bc.haveBlock(h)
	bc.RUnlock()
	return hasBlock
}

// GenesisBlock returns first basic block of the blockchain
func (bc *BlockChain) GenesisBlock() (*Block, error) {
	if bc.genesisBlock == nil {
		bc.genesisBlock = &Block{}
		genesisTransactionBytes, err := hex.DecodeString(bc.Network.GenesisCoinbaseTxHex)
		reader := bytes.NewReader(genesisTransactionBytes)

		if err != nil {
			return nil, err
		}

		if err := bc.genesisBlock.CoinbaseTransaction.Deserialize(reader); err != nil {
			return nil, err
		}

		bc.genesisBlock.MajorVersion = config.BlockMajorVersion1
		bc.genesisBlock.MinorVersion = config.BlockMinorVersion0
		bc.genesisBlock.Timestamp = bc.Network.GenesisTimestamp
		bc.genesisBlock.Nonce = bc.Network.GenesisNonce
	}

	return bc.genesisBlock, nil
}

// deserializeTransactions deserializes transactions to object, transactions are passing basic data validation
func (bc *BlockChain) deserializeTransactions(blogger *log.Entry, rt [][]byte) ([]Transaction, uint64, error) {
	transactions := make([]Transaction, len(rt))
	transactionsSize := uint64(0)

	for i, t := range transactions {
		tsSize := uint64(len(rt[i]))
		tsLogger := blogger.WithFields(log.Fields{
			"transaction_size":  tsSize,
			"transaction_index": i,
		})

		if tsSize > bc.Network.MaxTxSize {
			err := ErrAddBlockTransactionSizeMax
			tsLogger.Error(err)
			return nil, 0, err
		}

		r := bytes.NewReader(rt[i])
		if err := t.Deserialize(r); err != nil {
			tsLogger.Error(fmt.Errorf("%s: %w", ErrAddBlockTransactionDeserialization.Error(), err))
			return nil, 0, ErrAddBlockTransactionDeserialization
		}

		transactionsSize += tsSize
	}
	return transactions, transactionsSize, nil
}
