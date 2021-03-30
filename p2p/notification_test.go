package p2p

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDecodeResponseGetObject(t *testing.T) {
	payload, err := ioutil.ReadFile("./fixtures/2004.dat")
	assert.Nil(t, err)

	var rsp NotificationResponseGetObjects
	err = binary.Unmarshal(payload, &rsp)
	assert.Nil(t, err)

	assert.Equal(t, 588024, int(rsp.CurrentBlockchainHeight))

	assert.Len(t, rsp.Blocks, 128)
	assert.Len(t, rsp.Blocks[0].Block, 355)
	assert.Len(t, rsp.Blocks[1].Block, 355)
	assert.Len(t, rsp.Blocks[2].Block, 355)
	assert.Len(t, rsp.Blocks[0].Transactions, 0)
	assert.Len(t, rsp.Blocks[1].Transactions, 0)
	assert.Len(t, rsp.Blocks[2].Transactions, 0)

	assert.Len(t, rsp.MissedIds, 0)

	enc, err := binary.Marshal(rsp)
	assert.Nil(t, err)

	assert.Equal(t, payload, enc)

	var dec NotificationResponseGetObjects
	err = binary.Unmarshal(enc, &dec)
	assert.Nil(t, err)

	assert.Equal(t, rsp, dec)
}

func TestDecodeNewLiteObject(t *testing.T) {
	payload, err := ioutil.ReadFile("./fixtures/2009.dat")
	assert.Nil(t, err)

	var n NotificationNewLiteBlock
	err = binary.Unmarshal(payload, &n)
	assert.Nil(t, err)

	assert.Len(t, n.Block, 617)
	assert.Equal(t, 588158, int(n.CurrentBlockchainHeight))
	assert.Equal(t, 2, int(n.Hop))

	enc, err := binary.Marshal(n)
	assert.Nil(t, err)

	var dec NotificationNewLiteBlock
	err = binary.Unmarshal(enc, &dec)
	assert.Nil(t, err)

	assert.Equal(t, n, dec)
}

func TestRawBlock_ToBlock20(t *testing.T) {
	blockPayload, err := ioutil.ReadFile("./fixtures/block_20.dat")
	assert.Nil(t, err)

	transPayload1, err := ioutil.ReadFile("./fixtures/block_20_trans_0.dat")
	assert.Nil(t, err)

	rb := &RawBlock{
		Block: blockPayload,
		Transactions: [][]byte{transPayload1},
	}

	block, err := rb.ToBlock()
	if err != nil {
		panic(err)
	}

	hash, err := block.Hash()
	assert.Nil(t, err)

	assert.Equal(t, config.BlockMinorVersion0, block.BlockHeader.MinorVersion)
	assert.Equal(t, config.BlockMajorVersion1, block.BlockHeader.MajorVersion)
	assert.Equal(t, uint64(1464598015), block.BlockHeader.Timestamp)
	assert.Equal(t, uint32(2585362670), block.BlockHeader.Nonce)

	assert.Equal(t,
		"fdccbe0ae9966d138716669c8bf49dc2134a538b34dfb35c70e41f12e6605ec8",
		hash.String(),
	)

	assert.Equal(t,
		"1e5f81d4404082badcf90ee4f8301adbfc4cbe6dc78dfb864bcd8ba1b8b14b52",
		block.BlockHeader.Prev.String(),
	)

	assert.Equal(t, config.TransactionVersion1, block.Transaction.Version)
	assert.Equal(t, uint64(31), block.Transaction.UnlockHeight)
	assert.Len(t, block.Transaction.Inputs, 1)
	assert.Len(t, block.Transaction.Outputs, 6)
	assert.Equal(t, []byte{0x1, 0x16, 0x8d, 0xbe, 0x0, 0x87, 0x71, 0x64, 0xa7, 0x33, 0xf, 0x18, 0x3c, 0x3d, 0xbf, 0x53, 0xce, 0x4, 0x21, 0xe8, 0x3a, 0xe9, 0x74, 0x76, 0xfd, 0x2, 0x55, 0xbc, 0x4a, 0x74, 0xcd, 0xe0, 0x59}, block.Transaction.Extra)
}

func TestRawBlock_ToBlock123(t *testing.T) {
	bp, err := ioutil.ReadFile("./fixtures/block_123.dat")
	if err != nil {
		panic(err)
	}

	t1, err := ioutil.ReadFile("./fixtures/block_123_trans_0.dat")
	if err != nil {
		panic(err)
	}

	t2, err := ioutil.ReadFile("./fixtures/block_123_trans_1.dat")
	if err != nil {
		panic(err)
	}

	rb := &RawBlock{
		Block: bp,
		Transactions: [][]byte{t1, t2},
	}

	block, err := rb.ToBlock()
	if err != nil {
		panic(err)
	}

	hash, err := block.Hash()
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"40cd79fbbd5d0e86029255d2ad7e0410e3dae57f4758052dc8a44d72b0a5a436",
		hash.String(),
	)

	assert.Equal(t,
		"e5d41f68fd1f40f9698c95a95c72edc5376e449701f1c158dd08b93bc5461255",
		block.BlockHeader.Prev.String(),
	)
}
