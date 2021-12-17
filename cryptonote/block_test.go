package cryptonote

import (
	"bytes"
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

	hash := block.Hash()

	assert.Equal(t, &genesisBlockHash, hash)
}

func TestBlock_Deserialize200054(t *testing.T) {
	payload, _ := ioutil.ReadFile("./fixtures/block_200054.dat")

	var block Block
	r := bytes.NewReader(payload)
	err := block.Deserialize(r)
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"6769241077017f26c0a170fd9630c292695039399b6a22edaf293b52f2d542fb",
		block.PreviousBlockHash.String(),
	)

	hash := block.Hash()

	assert.Equal(t,
		"231a4584e0c13325024059482fabd99188574f51336d19c0b5787f7ccc9e4dfc",
		hash.String(),
	)

	thash := block.CoinbaseTransaction.Hash()

	assert.Equal(t,
		"25fc20b292ace0458ed2f9cf046588f3e7ecc02b6c2f372a29e4475545f9cef1",
		thash.String(),
	)

	b := block.Serialize()

	assert.Equal(t, payload, b)
}

func TestBlock_Deserialize105385(t *testing.T) {
	payload, _ := ioutil.ReadFile("./fixtures/block_105385.dat")

	var block Block
	r := bytes.NewReader(payload)
	err := block.Deserialize(r)
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"cc20ae5bd6c75e25a0885bcbb058e31c5b344dedc43a7f50c8ac6f1eaada795f",
		block.PreviousBlockHash.String(),
	)

	hash := block.Hash()
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"b8b793a00e0a1bb790987e5f6a1b551f9e397be5aa74595335094063be31f878",
		hash.String(),
	)

	thash := block.CoinbaseTransaction.Hash()

	assert.Equal(t,
		"effd9d3ffc37dc92f1f9721858c8ec11333c74ab74dcb60fc32926aa0dd51d8b",
		thash.String(),
	)

	b := block.Serialize()

	assert.Equal(t, payload, b)
}

func TestBlock_Deserialize60000(t *testing.T) {
	payload, _ := ioutil.ReadFile("./fixtures/block_60001.dat")

	var block Block
	r := bytes.NewReader(payload)
	err := block.Deserialize(r)
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"4cab277ce1d96569e6ec406c589f08468a490aafd729fccae3b46c7ba4cce3a7",
		block.Parent.Prev.String(),
	)

	assert.Equal(t,
		"4cab277ce1d96569e6ec406c589f08468a490aafd729fccae3b46c7ba4cce3a7",
		block.PreviousBlockHash.String(),
	)

	hash := block.Hash()

	assert.Equal(t,
		"8e39967eb50b8a922cbfe22fe02989218345cbd61ae651ddbecf00834910ff50",
		hash.String(),
	)

	thash := block.CoinbaseTransaction.Hash()

	assert.Equal(t,
		"20c2e650f05d3271a5064af2dcdada7ff1f79d8791dfb390b3f0a010942aba39",
		thash.String(),
	)

	b := block.Serialize()

	assert.Equal(t, payload, b)
}

func TestBlock_Deserialize(t *testing.T) {
	payload1, err := ioutil.ReadFile("./fixtures/block1.dat")
	assert.Nil(t, err)

	var block1 Block
	err = block1.Deserialize(bytes.NewReader(payload1))
	assert.Nil(t, err)

	hash1 := block1.Hash()

	assert.Equal(t, config.BlockMinorVersion0, block1.BlockHeader.MinorVersion)
	assert.Equal(t, config.BlockMajorVersion1, block1.BlockHeader.MajorVersion)
	assert.Equal(t, uint64(1464595534), block1.BlockHeader.Timestamp)
	assert.Equal(t, uint32(769685647), block1.BlockHeader.Nonce)

	assert.Equal(t,
		"93fd06c51fd8a6fc9db100adbdb4c1de11270a5186b790b454db8a7419c5615e",
		hash1.String(),
	)

	assert.Equal(t,
		"3125b79e4a42f8d4d2fc4dffea8442e185ebda940ecd4d3b449056a4ea0efea4",
		block1.BlockHeader.PreviousBlockHash.String(),
	)

	assert.Equal(t, config.TransactionVersion1, block1.CoinbaseTransaction.Version)
	assert.Equal(t, uint64(11), block1.CoinbaseTransaction.UnlockHeight)
	assert.Len(t, block1.CoinbaseTransaction.Inputs, 1)
	assert.Len(t, block1.CoinbaseTransaction.Outputs, 7)
	assert.Equal(t, []byte{0x1, 0xc1, 0xc9, 0xaa, 0xd0, 0xaa, 0x73, 0xbb, 0x5f, 0xc, 0x8, 0xb6, 0xb0, 0xe6, 0xe1, 0x4e, 0xc0, 0xdd, 0xa, 0xca, 0xa5, 0x6b, 0x9c, 0x52, 0x85, 0x74, 0xbd, 0x39, 0x29, 0x1c, 0xb4, 0x84, 0xc0}, block1.CoinbaseTransaction.Extra)

	payload2, err := ioutil.ReadFile("./fixtures/block2.dat")
	assert.Nil(t, err)

	var block2 Block
	err = block2.Deserialize(bytes.NewReader(payload2))
	assert.Nil(t, err)

	assert.Equal(t, config.BlockMinorVersion0, block2.BlockHeader.MinorVersion)
	assert.Equal(t, config.BlockMajorVersion1, block2.BlockHeader.MajorVersion)
	assert.Equal(t, uint32(1748233149), block2.BlockHeader.Nonce)
	assert.Equal(t, uint64(1464595535), block2.BlockHeader.Timestamp)

	assert.Equal(t, hash1.String(), block2.PreviousBlockHash.String())

	assert.Equal(t, config.TransactionVersion1, block2.CoinbaseTransaction.Version)
	assert.Equal(t, uint64(12), block2.CoinbaseTransaction.UnlockHeight)
	assert.Len(t, block2.CoinbaseTransaction.Inputs, 1)
	assert.Len(t, block2.CoinbaseTransaction.Outputs, 7)
	assert.Equal(t, []byte{0x1, 0x6f, 0x7f, 0x61, 0xe2, 0x4e, 0xfe, 0x12, 0x41, 0xc2, 0x55, 0xc8, 0x8, 0xc0, 0x95, 0xbb, 0x3a, 0x80, 0xd5, 0x93, 0x28, 0x1, 0x3d, 0xb0, 0x93, 0x55, 0x91, 0xaf, 0xf5, 0x5d, 0xf4, 0x55, 0xf1}, block2.CoinbaseTransaction.Extra)
}
