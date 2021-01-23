package cryptonote

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
	"io"

	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

// Key is any type of key
type Key interface {
	// Hex representation of the key
	Hex() string

	// Byte slice
	Bytes() *[32]byte
}

type key struct {
	b [32]byte // private key bytes
}

// Hex representation of key
func (k *key) Hex() string {
	return hex.EncodeToString(k.b[:])
}

// Bytes represention of key
func (k *key) Bytes() *[32]byte {
	return &k.b
}

// KeyFromHex returns key from hex string
func KeyFromHex(s string) (Key, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return KeyFromBytes(&decoded)
}

// KeyFromBytes key from bytes
func KeyFromBytes(b *[]byte) (Key, error) {
	if len(*b) != 32 {
		return nil, errors.New("Key must be 32 bytes length")
	}

	var keyBytes [32]byte
	copy(keyBytes[:], *b)
	return KeyFromArray(&keyBytes), nil
}

// KeyFromArray generate from array
func KeyFromArray(b *[32]byte) Key {
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
	return KeyFromArray(&keyBytes)
}

// ViewFromSpend returns deterministic private key
func ViewFromSpend(k Key) Key {
	kbytes := k.Bytes()[:]
	khash := hash.Keccak(&kbytes)
	key := KeyFromArray(reduceBytesToPoint(&khash))

	return key
}

// GenerateKey cryptonote keys
func GenerateKey() (Key, error) {
	randomBytes := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return nil, err
	}

	return KeyFromArray(reduceBytesToPoint(&randomBytes)), nil
}

func reduceBytesToPoint(b *[]byte) *[32]byte {
	var in [64]byte
	var out [32]byte

	copy(in[:], *b)
	ed.ScReduce(&out, &in)
	return &out
}
