package crypto

import (
	"crypto/rand"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"io"
)

type EllipticCurvePoint [32]byte

type EllipticCurveScalar [32]byte

// Check that the point is on curve
func (p *EllipticCurvePoint) Check() bool {
	var point ed.ExtendedGroupElement
	return point.FromBytes((*[32]byte)(p))
}

// RandomScalar generates random scalar
func RandomScalar() EllipticCurveScalar {
	var randomBytes [64]byte

	if _, err := io.ReadFull(rand.Reader, randomBytes[:]); err != nil {
		panic(err)
	}

	return ed.ScReduce(randomBytes)
}
