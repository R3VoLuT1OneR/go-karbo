package crypto

import (
	"encoding/hex"
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

var (
	ErrKeyNotOnCurve  = errors.New("key is not on curve")
	ErrKeyWrongLength = errors.New("key must be 32 bytes length")
)

// Key is any type of key
type Key interface {
	// Hex representation of the key
	Hex() string

	// Bytes array
	Bytes() [32]byte

	// BytesSlice return bytes slice
	BytesSlice() []byte

	// Check if key is it is valid key
	Check() bool
}

type key struct {
	b [32]byte // private key bytes
}

// Hex representation of key
func (k *key) Hex() string {
	return hex.EncodeToString(k.b[:])
}

// Bytes representation of the key
func (k *key) Bytes() [32]byte {
	return k.b
}

func (k *key) BytesSlice() []byte {
	return k.b[:]
}

// Check that the point is on curve
func (k *key) Check() bool {
	var point ed.ExtendedGroupElement
	b := k.Bytes()
	return point.FromBytes(&b)
}

// KeyFromHex returns key from hex string
func KeyFromHex(s string) (Key, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return KeyFromBytes(decoded)
}

// KeyFromBytes key from bytes
func KeyFromBytes(b []byte) (Key, error) {
	if len(b) != 32 {
		return nil, ErrKeyWrongLength
	}

	var keyBytes [32]byte
	copy(keyBytes[:], b)
	return KeyFromArray(keyBytes), nil
}

// KeyFromArray generate from array
func KeyFromArray(b [32]byte) Key {
	return &key{b}
}

// PublicFromPrivate key from private
func PublicFromPrivate(privateKey Key) (Key, error) {
	if !ed.ScCheck(privateKey.Bytes()) {
		return nil, ErrKeyNotOnCurve
	}

	var point ed.ExtendedGroupElement
	ed.GeScalarMultBase(&point, privateKey.Bytes())

	return KeyFromArray(point.ToBytes()), nil
}

// ViewFromSpend returns deterministic private key
func ViewFromSpend(k Key) Key {
	var reduceBytes [64]byte
	copy(reduceBytes[:], Keccak(k.BytesSlice()))
	key := KeyFromArray(ed.ScReduce(reduceBytes))

	return key
}

// GenerateKey cryptonote keys
func GenerateKey() (Key, error) {
	randomScalar := RandomScalar()

	return KeyFromArray(randomScalar), nil
}