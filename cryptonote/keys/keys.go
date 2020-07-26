package keys

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
)

// Key is any type of key
type Key interface {
	// Hex represention of the key
	Hex() string

	// Byte slice
	Bytes() *[32]byte
}

type key struct {
	b [32]byte // private key bytes
}

// Hex represention of key
func (k *key) Hex() string {
	return hex.EncodeToString(k.b[:])
}

// Bytes represention of key
func (k *key) Bytes() *[32]byte {
	return &k.b
}

// FromHex returns key from hex string
func FromHex(s string) (Key, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return FromBytes(&decoded)
}

// FromBytes key from bytes
func FromBytes(b *[]byte) (Key, error) {
	if len(*b) != 32 {
		return nil, errors.New("Key must be 32 bytes length")
	}

	var keyBytes [32]byte
	copy(keyBytes[:], *b)
	return FromArray(&keyBytes), nil
}

// FromArray generate from array
func FromArray(b *[32]byte) Key {
	return &key{*b}
}

// PublicFromPrivate key from private
func PublicFromPrivate(k Key) Key {
	if !ed.ScCheck(k.Bytes()) {
		panic(errors.New("Provided key is not on curve"))
	}

	var point ed.ExtendedGroupElement
	ed.GeScalarMultBase(&point, k.Bytes())

	var keyBytes [32]byte
	point.ToBytes(&keyBytes)
	return FromArray(&keyBytes)
}

// ViewFromSpend returns deterministic private key
func ViewFromSpend(k Key) Key {
	khash := hash.Keccak(k.Bytes()[:])
	key := FromArray(reduceBytesToPoint(&khash))

	return key
}

// GenerateKey cryptonote keys
func GenerateKey() (Key, error) {
	randomBytes := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return nil, err
	}

	return FromArray(reduceBytesToPoint(&randomBytes)), nil
}

func reduceBytesToPoint(b *[]byte) *[32]byte {
	var in [64]byte
	var out [32]byte

	copy(in[:], *b)
	ed.ScReduce(&out, &in)
	return &out
}
