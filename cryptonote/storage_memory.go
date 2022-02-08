package cryptonote

import "github.com/r3volut1oner/go-karbo/crypto"

type memoryStorage struct {

	// index keeps the block height for specific height
	index map[crypto.Hash]uint32

	heightIndex map[uint32]*crypto.Hash
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		index:       map[crypto.Hash]uint32{},
		heightIndex: map[uint32]*crypto.Hash{},
	}
}

func (s *memoryStorage) GetBlockHashByHeight(height uint32) (*crypto.Hash, error) {
	return s.heightIndex[height], nil
}
