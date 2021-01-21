package address

import (
	"errors"

	"github.com/r3volut1oner/go-karbo/cryptonote/base58"
	"github.com/r3volut1oner/go-karbo/cryptonote/keys"
)

// Address represents cryptonote address
type Address struct {
	Tag            uint64
	SpendPublicKey keys.Key
	ViewPublicKey  keys.Key
	Base58         string
}

// FromString provides address from string
func FromString(s string) (a Address, err error) {
	tag, data, err := base58.DecodeAddr(s)

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
	a.SpendPublicKey = keys.FromArray(&spendPublicKeyData)
	a.ViewPublicKey = keys.FromArray(&viewPublicKeyData)
	a.Tag = tag

	return
}

// Generate from tag and public keys
func Generate(tag uint64, sk, vk keys.Key) (a Address) {
	a.Tag = tag
	a.SpendPublicKey = sk
	a.ViewPublicKey = vk

	var b []byte
	b = append(b, sk.Bytes()[:]...)
	b = append(b, vk.Bytes()[:]...)

	a.Base58 = base58.EncodeAddr(tag, b)

	return
}
