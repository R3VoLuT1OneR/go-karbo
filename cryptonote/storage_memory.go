package cryptonote

// Memory storage keeps all the blockchain in the memory.
// This is the basic implementation of the Storage interface.

import (
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/utils"
	"sync"
)

type memoryStorage struct {

	// blockIndex keeps block hash for specific index
	blockIndex map[uint32]*Block

	// blockInfosIndex keeps block infos for specific index
	blockInfosIndex map[uint32]*blockInfo

	blockInfosHashIndex map[crypto.Hash]*blockInfo

	transactionsIndex map[uint32]*[]Transaction

	spentKeysImagesIndex map[uint32]*[]crypto.KeyImage

	spentMultisignatureGlobalIndexesIndex map[uint32]*[]MultisigAmountGlobalOutputIndexPair

	topBlock *Block

	sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		blockIndex:                            map[uint32]*Block{},
		blockInfosIndex:                       map[uint32]*blockInfo{},
		transactionsIndex:                     map[uint32]*[]Transaction{},
		spentKeysImagesIndex:                  map[uint32]*[]crypto.KeyImage{},
		spentMultisignatureGlobalIndexesIndex: map[uint32]*[]MultisigAmountGlobalOutputIndexPair{},
	}
}

func (s *memoryStorage) Init(genesisBlock *Block) error {
	info := blockInfo{
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

func (s *memoryStorage) PushBlock(block *Block, info *blockInfo, details TransactionsDetails) error {
	s.Lock()
	err := s.appendBlock(block, info, details)
	s.Unlock()
	return err
}

func (s *memoryStorage) appendBlock(block *Block, info *blockInfo, details TransactionsDetails) error {
	hash := block.Hash()
	index := block.Index()

	if block.Index() != info.Index {
		return utils.AssertionError("block info index and block index must be same")
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

	s.blockIndex[index] = block
	s.blockInfosIndex[index] = info
	s.blockInfosHashIndex[*hash] = info

	s.transactionsIndex[index] = &details.transactions
	s.spentKeysImagesIndex[index] = &details.spentKeyImages
	s.spentMultisignatureGlobalIndexesIndex[index] = &details.spentMultisignatureGlobalIndexes

	if s.topBlock == nil || index > s.topBlock.Index() {
		s.topBlock = block
	}

	return nil
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

func (s *memoryStorage) getBlockInfoAtIndex(index uint32) *blockInfo {
	s.RLock()
	info := s.blockInfosIndex[index]
	s.RUnlock()
	return info
}

func (s *memoryStorage) Close() error {
	return nil
}
