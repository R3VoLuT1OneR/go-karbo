package cryptonote

import (
	"bytes"
	"encoding/hex"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_Deserialize(t *testing.T) {
	network := config.MainNet()

	serializedTransaction, err := hex.DecodeString(network.GenesisCoinbaseTxHex)
	if err != nil {
		panic(err)
	}

	var transaction Transaction
	reader := bytes.NewReader(serializedTransaction)
	err = transaction.Deserialize(reader)
	assert.Nil(t, err)

	assert.Equal(t, config.TransactionVersion1, transaction.TransactionPrefix.Version)
	assert.Equal(t, 10, int(transaction.TransactionPrefix.UnlockHeight))

	assert.Len(t, transaction.TransactionPrefix.Inputs, 1)
	assert.IsType(t, InputCoinbase{}, transaction.TransactionPrefix.Inputs[0])
	assert.IsType(t, uint32(1), transaction.TransactionPrefix.Inputs[0].(InputCoinbase).Height)

	assert.Len(t, transaction.TransactionPrefix.Outputs, 1)
	assert.IsType(t, TransactionOutput{}, transaction.TransactionPrefix.Outputs[0])
	assert.Equal(t, 38146972656250, int(transaction.TransactionPrefix.Outputs[0].Amount))
	assert.IsType(t, OutputKey{}, transaction.TransactionPrefix.Outputs[0].Target)
	assert.Equal(t,
		[32]byte{
			0x9b, 0x2e, 0x4c, 0x2, 0x81, 0xc0, 0xb0, 0x2e,
			0x7c, 0x53, 0x29, 0x1a, 0x94, 0xd1, 0xd0, 0xcb,
			0xff, 0x88, 0x83, 0xf8, 0x2, 0x4f, 0x51, 0x42,
			0xee, 0x49, 0x4f, 0xfb, 0xbd, 0x8, 0x80, 0x71,
		},
		*transaction.TransactionPrefix.Outputs[0].Target.(OutputKey).Key.Bytes(),
	)

	assert.Equal(t,
		[]byte{
			0x1, 0xf9, 0x4, 0x92, 0x5c, 0xc2, 0x3f, 0x86,
			0xf9, 0xf3, 0x56, 0x51, 0x88, 0x86, 0x22, 0x75,
			0xdc, 0x55, 0x6a, 0x9b, 0xdf, 0xb6, 0xae, 0xc2,
			0x2c, 0x5a, 0xca, 0x7f, 0x1, 0x77, 0xc4, 0x5b, 0xa8,
		},
		transaction.TransactionPrefix.Extra,
	)

	serialized, err := transaction.Serialize()
	assert.Nil(t, err)
	assert.Equal(t, serializedTransaction, serialized)

	expectedHash := Hash{
		0x11, 0xa, 0xf2, 0xe4, 0x2d, 0xd6, 0x29, 0xe3,
		0xda, 0x49, 0xec, 0xe8, 0xc, 0x40, 0x7, 0xc9,
		0xe, 0xc3, 0x20, 0xa8, 0xa4, 0x55, 0xcf, 0xd2,
		0x2c, 0x9b, 0x80, 0x6a, 0x78, 0x69, 0xa4, 0x15,
	}

	hash, err := transaction.Hash()

	assert.Nil(t, err)
	assert.Equal(t, expectedHash, *hash)
}
