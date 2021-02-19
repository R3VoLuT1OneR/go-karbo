package p2p

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEncodeRequestChain(t *testing.T) {
	rc, err := newRequestChain(config.MainNet())
	assert.Nil(t, err)

	b, err := binary.Marshal(*rc)
	assert.Nil(t, err)

	var d NotificationRequestChain
	err = binary.Unmarshal(b, &d)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)

	assert.Equal(t, *rc, d)
}

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
