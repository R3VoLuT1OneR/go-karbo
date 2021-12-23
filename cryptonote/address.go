package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto"
)

// Address represents cryptonote address
type Address struct {
	Tag            uint64
	SpendPublicKey crypto.Key
	ViewPublicKey  crypto.Key
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

	var spendPublicKeyData [32]byte
	var viewPublicKeyData [32]byte

	copy(spendPublicKeyData[:], data[:32])
	copy(viewPublicKeyData[:], data[32:64])

	a.Base58 = s
	a.SpendPublicKey = crypto.KeyFromArray(spendPublicKeyData)
	a.ViewPublicKey = crypto.KeyFromArray(viewPublicKeyData)
	a.Tag = tag

	return
}

// Generate from tag and public keys
func Generate(tag uint64, sk, vk crypto.Key) (a Address) {
	a.Tag = tag
	a.SpendPublicKey = sk
	a.ViewPublicKey = vk

	var b []byte
	b = append(b, sk.BytesSlice()[:]...)
	b = append(b, vk.BytesSlice()[:]...)

	a.Base58 = EncodeAddr(tag, b)

	return
}
