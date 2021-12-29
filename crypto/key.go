package crypto

import (
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

type SecretKey EllipticCurveScalar

type PublicKey EllipticCurvePoint

var (
	ErrKeyNotOnCurve  = errors.New("key is not on curve")
	ErrKeyWrongLength = errors.New("key must be 32 bytes length")
)

// Check that public key is on curve
func (publicKey *PublicKey) Check() bool {
	ecPoint := EllipticCurvePoint(*publicKey)
	return ecPoint.Check()
}

// PublicFromSecret key from private
func PublicFromSecret(secretKey *SecretKey) (*PublicKey, error) {
	if !ed.ScCheck(*secretKey) {
		return nil, ErrKeyNotOnCurve
	}

	var point ed.ExtendedGroupElement
	ed.GeScalarMultBase(&point, *secretKey)

	pk := PublicKey(point.ToBytes())
	return &pk, nil
}

// ViewFromSpend returns deterministic private key
func ViewFromSpend(spendKey *SecretKey) (viewKey SecretKey) {
	var reduceBytes [64]byte
	copy(reduceBytes[:], Keccak(spendKey[:]))
	viewKey = ed.ScReduce(reduceBytes)

	return viewKey
}

// GenerateKey cryptonote keys
func GenerateKey() (SecretKey, error) {
	return SecretKey(RandomScalar()), nil
}
