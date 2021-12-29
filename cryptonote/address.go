package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto"
)

// Address represents cryptonote address
type Address struct {
	Tag            uint64
	SpendPublicKey crypto.PublicKey
	ViewPublicKey  crypto.PublicKey
	Base58         string
}

// FromString provides address from string
func FromString(s string) (a Address, err error) {
	tag, data, err := DecodeAddr(s)

	if err != nil {
		return
	}

	if len(data) != 64 {
		err = errors.New("encoded data has wrong length")

		return
	}

	var spendPublicKeyBytes [32]byte
	var viewPublicKeyBytes [32]byte

	copy(spendPublicKeyBytes[:], data[:32])
	copy(viewPublicKeyBytes[:], data[32:64])

	a.Base58 = s
	a.SpendPublicKey = spendPublicKeyBytes
	a.ViewPublicKey = viewPublicKeyBytes
	a.Tag = tag

	return
}

// Generate from tag and public keys
func Generate(tag uint64, spendPublicKey, viewPublicKey crypto.PublicKey) (a Address) {
	a.Tag = tag
	a.SpendPublicKey = spendPublicKey
	a.ViewPublicKey = viewPublicKey

	var b []byte
	b = append(b, spendPublicKey[:]...)
	b = append(b, viewPublicKey[:]...)

	a.Base58 = EncodeAddr(tag, b)

	return
}
