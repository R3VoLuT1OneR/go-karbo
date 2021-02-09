package p2p

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/stretchr/testify/assert"
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