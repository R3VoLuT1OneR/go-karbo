package keys

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"github.com/r3volut1oner/go-karbo/cryptonote/hash"
)

// PrivateKey represents cryptonote private key
type PrivateKey struct {
	b [32]byte   // private key bytes
	p *PublicKey // public key cache
}

// PublicKey represents cryptonote public key
type PublicKey struct {
	b [32]byte
}

// Key is any type of key
type Key interface {
	Hex() string
}

// Public key from private
func (k *PrivateKey) Public() *PublicKey {
	if k.p != nil {
		var private [32]byte = k.b
		var public [32]byte

		var point ed.ExtendedGroupElement
		ed.GeScalarMultBase(&point, &private)
		point.ToBytes(&public)

		k.p = &PublicKey{public}
	}

	return k.p
}

// Hex represention of key
func (k *PrivateKey) Hex() string {
	return hex.EncodeToString(k.b[:])
}

// Hex represention of key
func (k *PublicKey) Hex() string {
	return hex.EncodeToString(k.b[:])
}

// PrivateFromHex private key from string
func PrivateFromHex(s string) (p PrivateKey, err error) {
	p.b, err = keyFromHex(s)

	if err != nil {
		return
	}

	return
}

// GenerateViewFromSpend returns deterministic private key
func GenerateViewFromSpend(k *PrivateKey) PrivateKey {
	var khash = hash.Keccak(k.b[:32])

	var pkey PrivateKey
	copy(pkey.b[:], khash[:32])

	return GenerateDeterministicKey(&pkey)
}

// GenerateDeterministicKey compute determenistic key from private key
func GenerateDeterministicKey(k *PrivateKey) (new PrivateKey) {
	var in [64]byte
	copy(in[:], k.b[:])
	ed.ScReduce(&new.b, &in)

	return
}

// GenerateSpendKey cryptonote keys
func GenerateSpendKey() (k PrivateKey, err error) {
	randomBytes := make([]byte, 64)
	if _, err = io.ReadFull(rand.Reader, randomBytes); err != nil {
		return
	}

	var seed [64]byte
	var key [32]byte

	copy(seed[:], randomBytes[:64])
	ed.ScReduce(&key, &seed)

	k.b = key

	return
}

func keyFromHex(s string) (k [32]byte, err error) {
	var bkey []byte

	bkey, err = hex.DecodeString(s)
	if err != nil {
		return
	}

	copy(k[:], bkey[:32])

	return
}
