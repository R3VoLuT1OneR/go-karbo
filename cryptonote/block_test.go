package cryptonote

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestGenerateGenesisBlock(t *testing.T) {
	net := config.MainNet()

	genesisBlockHash := Hash{
		0x31, 0x25, 0xb7, 0x9e, 0x4a, 0x42, 0xf8, 0xd4,
		0xd2, 0xfc, 0x4d, 0xff, 0xea, 0x84, 0x42, 0xe1,
		0x85, 0xeb, 0xda, 0x94, 0xe, 0xcd, 0x4d, 0x3b,
		0x44, 0x90, 0x56, 0xa4, 0xea, 0xe, 0xfe, 0xa4,
	}

	block, err := GenerateGenesisBlock(net)

	// TODO: Add tests for genesis block
	assert.Nil(t, err)

	hash, err := block.Hash()

	assert.Nil(t, err)
	assert.Equal(t, &genesisBlockHash, hash)
}

func TestBlock_Deserialize(t *testing.T) {
	payload1, err := ioutil.ReadFile("./fixtures/block1.dat")
	assert.Nil(t, err)

	var block1 Block
	err = block1.Deserialize(payload1)
	assert.Nil(t, err)

	assert.Equal(t, config.BlockMinorVersion0, block1.BlockHeader.MinorVersion)
	assert.Equal(t, config.BlockMajorVersion1, block1.BlockHeader.MajorVersion)
	assert.Equal(t, uint64(1464595534), block1.BlockHeader.Timestamp)
	assert.Equal(t, uint32(769685647), block1.BlockHeader.Nonce)

	assert.Equal(t,
		"3125b79e4a42f8d4d2fc4dffea8442e185ebda940ecd4d3b449056a4ea0efea4",
		block1.BlockHeader.Prev.String(),
	)

	assert.Equal(t, config.TransactionVersion1, block1.Transaction.Version)
	assert.Equal(t, uint64(11), block1.Transaction.UnlockHeight)
	assert.Len(t, block1.Transaction.Inputs, 1)
	assert.Len(t, block1.Transaction.Outputs, 7)
	assert.Equal(t, []byte{0x1, 0xc1, 0xc9, 0xaa, 0xd0, 0xaa, 0x73, 0xbb, 0x5f, 0xc, 0x8, 0xb6, 0xb0, 0xe6, 0xe1, 0x4e, 0xc0, 0xdd, 0xa, 0xca, 0xa5, 0x6b, 0x9c, 0x52, 0x85, 0x74, 0xbd, 0x39, 0x29, 0x1c, 0xb4, 0x84, 0xc0}, block1.Transaction.Extra)

	payload2, err := ioutil.ReadFile("./fixtures/block2.dat")
	assert.Nil(t, err)

	var block2 Block
	err = block2.Deserialize(payload2)
	assert.Nil(t, err)

	assert.Equal(t, config.BlockMinorVersion0, block2.BlockHeader.MinorVersion)
	assert.Equal(t, config.BlockMajorVersion1, block2.BlockHeader.MajorVersion)
	assert.Equal(t, uint32(1748233149), block2.BlockHeader.Nonce)
	assert.Equal(t, uint64(1464595535), block2.BlockHeader.Timestamp)

	hash, err := block1.Hash()
	assert.Nil(t, err)
	assert.Equal(t, hash.String(), block2.Prev.String())

	assert.Equal(t, config.TransactionVersion1, block2.Transaction.Version)
	assert.Equal(t, uint64(12), block2.Transaction.UnlockHeight)
	assert.Len(t, block2.Transaction.Inputs, 1)
	assert.Len(t, block2.Transaction.Outputs, 7)
	assert.Equal(t, []byte{0x1, 0x6f, 0x7f, 0x61, 0xe2, 0x4e, 0xfe, 0x12, 0x41, 0xc2, 0x55, 0xc8, 0x8, 0xc0, 0x95, 0xbb, 0x3a, 0x80, 0xd5, 0x93, 0x28, 0x1, 0x3d, 0xb0, 0x93, 0x55, 0x91, 0xaf, 0xf5, 0x5d, 0xf4, 0x55, 0xf1}, block2.Transaction.Extra)
}
