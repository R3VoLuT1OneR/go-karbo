package cryptonote

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockChain_NewBlockChain(t *testing.T) {
	network := config.TestNet()
	logger := logrus.New()

	bc := NewBlockChain(network, logger)

	assert.IsType(t, &BlockChain{}, bc)
	assert.Equal(t, network, bc.Network)
}

func TestBlockChain_GenesisBlock(t *testing.T) {
	net := config.MainNet()
	logger := logrus.New()

	genesisBlockHash := crypto.Hash{
		0x31, 0x25, 0xb7, 0x9e, 0x4a, 0x42, 0xf8, 0xd4,
		0xd2, 0xfc, 0x4d, 0xff, 0xea, 0x84, 0x42, 0xe1,
		0x85, 0xeb, 0xda, 0x94, 0xe, 0xcd, 0x4d, 0x3b,
		0x44, 0x90, 0x56, 0xa4, 0xea, 0xe, 0xfe, 0xa4,
	}

	bc := NewBlockChain(net, logger)
	block, err := bc.GenesisBlock()

	// TODO: Add tests for genesis block
	assert.Nil(t, err)

	hash := block.Hash()

	assert.Equal(t, &genesisBlockHash, hash)
	assert.Equal(t, block.PreviousBlockHash, crypto.Hash{})
}
