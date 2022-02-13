package cryptonote

// Memory storage keeps all the blockchain in the memory.
// This is the basic implementation of the Storage interface.

import (
	"github.com/r3volut1oner/go-karbo/crypto"
	"sync"
)

type memoryStorage struct {

	// hashIndex keeps the block index for specific hash
	hashIndex map[crypto.Hash]*Block

	// blockIndex keeps block hash for specific index
	blockIndex map[uint32]*Block

	// blockInfos keeps block infos for specific index
	blockInfos map[uint32]*blockInfo

	topBlock *Block

	sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		hashIndex:  map[crypto.Hash]*Block{},
		blockIndex: map[uint32]*Block{},
		blockInfos: map[uint32]*blockInfo{},
	}
}

func (s *memoryStorage) Init(genesisBlock *Block) error {
	info := blockInfo{
		Index:                0,
		CumulativeDifficulty: 1,
		Size:                 genesisBlock.BaseTransaction.Size(),
		TotalGeneratedCoins:  genesisBlock.BaseTransaction.Outputs[0].Amount,
	}

	err := s.PushBlock(genesisBlock, &info)
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

func (s *memoryStorage) PushBlock(block *Block, info *blockInfo) error {
	s.Lock()
	err := s.appendBlock(block, info)
	s.Unlock()
	return err
}

func (s *memoryStorage) appendBlock(block *Block, info *blockInfo) error {
	hash := block.Hash()
	index := block.Index()

	if s.haveBlock(hash) {
		return ErrStorageBlockExists
	}

	if _, ok := s.blockIndex[index]; ok {
		return ErrStorageBlockExists
	}

	s.hashIndex[*hash] = block
	s.blockIndex[index] = block
	s.blockInfos[index] = info

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
	block := s.hashIndex[*hash]
	s.RUnlock()
	return block
}

func (s *memoryStorage) haveBlock(hash *crypto.Hash) bool {
	if _, ok := s.hashIndex[*hash]; ok {
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
	info := s.blockInfos[index]
	s.RUnlock()
	return info
}

func (s *memoryStorage) Close() error {
	return nil
}
