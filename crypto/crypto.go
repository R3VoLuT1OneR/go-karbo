package crypto

import (
	"crypto/rand"
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"io"
)

type EllipticCurvePoint [32]byte

type EllipticCurveScalar [32]byte

type KeyImage EllipticCurvePoint

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

func GenerateKeyImage(publicKey *PublicKey, secretKey *SecretKey) (*KeyImage, error) {
	if !ed.ScCheck(*secretKey) {
		return nil, errors.New("wrong secret key provided")
	}

	pkHash := HashFromBytes(publicKey[:])
	point, err := pkHash.toEc()
	if err != nil {
		return nil, err
	}

	point2 := ed.GeScalarMult((*[32]byte)(secretKey), point)
	keyImage := KeyImage(point2.ToBytes())

	return &keyImage, nil
}
