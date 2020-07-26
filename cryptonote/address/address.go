package address

import (
	"errors"
	"fmt"

	"github.com/r3volut1oner/go-karbo/cryptonote/base58"
	"github.com/r3volut1oner/go-karbo/cryptonote/keys"
)

// Address represents cryptonote address
type Address struct {
	Tag            uint64
	SpendPublicKey keys.PublicKey
	ViewPublicKey  keys.PublicKey
	Base58         string
}

// FromString provides address from string
func FromString(s string) (a Address, err error) {
	tag, data, err := base58.DecodeAddr(s)

	if err != nil {
		return
	}

	if len(data) != 64 {
		err = errors.New("Encoded data has wrong length")

		return
	}

	fmt.Println("encoded", len(data[:32]), len(data[32:64]))

	a.Base58 = s
	a.SpendPublicKey, _ = keys.PublicFromBytes(data[:32])
	a.ViewPublicKey, _ = keys.PublicFromBytes(data[32:64])
	a.Tag = tag

	return
}

// Generate from tag and public keys
func Generate(tag uint64, sk, vk keys.PublicKey) (a Address) {
	a.Tag = tag
	a.SpendPublicKey = sk
	a.ViewPublicKey = vk

	b := append([]byte{}, sk.Bytes()...)
	b = append(b, vk.Bytes()...)
	a.Base58 = base58.EncodeAddr(tag, b)

	return
}
