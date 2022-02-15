package cryptonote

// Memory storage keeps all the blockchain in the memory.
// This is the basic implementation of the Storage interface.

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/utils"
	"math"
	"sync"
)

type memoryStorage struct {
	topBlock *Block

	// blockIndex keeps block hash for specific index
	blockIndex map[uint32]*Block

	// blockInfosIndex keeps block infos for specific index
	blockInfosIndex     map[uint32]*BlockInfo
	blockInfosHashIndex map[crypto.Hash]*BlockInfo

	spentKeysImagesIndex                  map[uint32]*[]crypto.KeyImage
	spentMultisignatureGlobalIndexesIndex map[uint32]*[]MultisigAmountGlobalOutputIndexPair

	transactionsIndex     map[uint32]*[]Transaction
	transactionInfosIndex map[uint32]*[]TransactionInfo

	sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		blockIndex: map[uint32]*Block{},

		blockInfosIndex:     map[uint32]*BlockInfo{},
		blockInfosHashIndex: map[crypto.Hash]*BlockInfo{},

		spentKeysImagesIndex:                  map[uint32]*[]crypto.KeyImage{},
		spentMultisignatureGlobalIndexesIndex: map[uint32]*[]MultisigAmountGlobalOutputIndexPair{},

		transactionsIndex: map[uint32]*[]Transaction{},
	}
}

func (s *memoryStorage) Init(genesisBlock *Block) error {
	info := BlockInfo{
		Index:                0,
		CumulativeDifficulty: 1,
		Size:                 genesisBlock.BaseTransaction.Size(),
		TotalGeneratedCoins:  genesisBlock.BaseTransaction.Outputs[0].Amount,
	}

	err := s.PushBlock(genesisBlock, &info, TransactionsDetails{})
	return err
}

func (s *memoryStorage) TopIndex() (uint32, error) {
	s.RLock()
	index := s.topBlock.Index()
	s.Unlock()
	return index, nil
}

func (s *memoryStorage) TopBlock() (*Block, error) {
	s.RLock()
	block := s.topBlock
	s.RUnlock()
	return block, nil
}

func (s *memoryStorage) PushBlock(block *Block, info *BlockInfo, details TransactionsDetails) error {
	s.Lock()
	err := s.appendBlock(block, info, details)
	s.Unlock()
	return err
}

func (s *memoryStorage) appendBlock(block *Block, info *BlockInfo, details TransactionsDetails) error {
	hash := block.Hash()
	index := block.Index()

	if block.Index() != info.Index {
		return utils.AssertionError("block info i and block i must be same")
	}

	if *block.Hash() != info.Hash {
		return utils.AssertionError("block info hash and block hash must be same")
	}

	if s.haveBlock(hash) {
		return ErrStorageBlockExists
	}

	if _, ok := s.blockIndex[index]; ok {
		return ErrStorageBlockExists
	}

	// Save all the block information
	s.blockIndex[index] = block
	s.blockInfosIndex[index] = info
	s.blockInfosHashIndex[*hash] = info

	// Save extra information about the block
	s.spentKeysImagesIndex[index] = &details.spentKeyImages
	s.spentMultisignatureGlobalIndexesIndex[index] = &details.spentMultisignatureGlobalIndexes

	// Save base transaction
	if err := s.pushTransaction(block, &block.BaseTransaction, 0); err != nil {
		return err
	}

	// Save base transactions
	for i := 1; i <= len(details.transactions); i++ {
		if err := s.pushTransaction(block, &details.transactions[i], uint32(i)); err != nil {
			return err
		}
	}

	if s.topBlock == nil || index > s.topBlock.Index() {
		s.topBlock = block
	}

	return nil
}

func (s *memoryStorage) pushTransaction(block *Block, tx *Transaction, txIndex uint32) error {
	transactionInfo := TransactionInfo{
		Index:      txIndex,
		BlockIndex: block.Index(),
		Hash:       *tx.Hash(),
		UnlockTime: tx.UnlockHeight,
		Outputs:    make([]TransactionOutputTarget, len(tx.Outputs)),
	}

	if len(tx.Outputs) > math.MaxUint16 {
		return utils.AssertionError("there are too much outputs in transactions")
	}

	keyIndexes := map[uint64][]PackedOutIndex{}
	multiIndexes := map[uint64][]PackedOutIndex{}

	for i, output := range tx.Outputs {
		transactionInfo.Outputs[i] = output.Target

		poi := PackedOutIndex{
			BlockIndex:       block.Index(),
			TransactionIndex: txIndex,
			Index:            uint32(i),
		}

		switch output.Target.(type) {
		case OutputKey:
			keyIndexes[output.Amount] = append(keyIndexes[output.Amount], poi)

		default:
			return errors.New("unknown output type")
		}
	}

	// TODO: Save global indexes
	// GlobalIndexes []uint32
	// AmountToKeyIndexes map[uint64][]uint32
	// AmountToMultiIndexes map[uint64][]uint32

}

func (s *memoryStorage) HaveBlock(hash *crypto.Hash) bool {
	s.RLock()
	have := s.haveBlock(hash)
	s.RUnlock()
	return have
}

func (s *memoryStorage) GetBlock(hash *crypto.Hash) *Block {
	s.RLock()
	defer s.RUnlock()

	if blockInfo, ok := s.blockInfosHashIndex[*hash]; ok {
		block := s.blockIndex[blockInfo.Index]
		return block
	}

	return nil
}

func (s *memoryStorage) haveBlock(hash *crypto.Hash) bool {
	if _, ok := s.blockInfosHashIndex[*hash]; ok {
		return true
	}

	return false
}

func (s *memoryStorage) HashAtIndex(index uint32) (*crypto.Hash, error) {
	s.RLock()
	defer s.RUnlock()

	if block, ok := s.blockIndex[index]; ok {
		return block.Hash(), nil
	}

	return nil, nil
}

func (s *memoryStorage) getBlockInfoAtIndex(index uint32) *BlockInfo {
	s.RLock()
	info := s.blockInfosIndex[index]
	s.RUnlock()
	return info
}

func (s *memoryStorage) Close() error {
	return nil
}
