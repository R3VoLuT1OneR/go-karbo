package crypto

import (
	"crypto/rand"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"io"
)

type EllipticCurvePoint [32]byte

type EllipticCurveScalar [32]byte

func HashToScalar(b []byte) EllipticCurveScalar {
	hashed := Keccak(b)
	var ba [32]byte
	copy(ba[:], hashed[:])
	return ed.ScReduce32(ba)
}

func RandomScalar() EllipticCurveScalar {
	var randomBytes [64]byte

	if _, err := io.ReadFull(rand.Reader, randomBytes[:]); err != nil {
		panic(err)
	}

	return ed.ScReduce(randomBytes)
}
