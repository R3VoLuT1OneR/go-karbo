package crypto

import (
	"crypto/rand"
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"io"
)

type EllipticCurvePoint [32]byte

type EllipticCurveScalar [32]byte

type SecretKey EllipticCurveScalar

type PublicKey EllipticCurvePoint

type KeyDerivation EllipticCurvePoint

type KeyImage EllipticCurvePoint

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

func GenerateKeyDerivation(publicKey PublicKey, secretKey SecretKey) (*KeyDerivation, error) {
	var point *ed.ExtendedGroupElement

	if !ed.ScCheck(secretKey) {
		return nil, errors.New("broken private key provided")
	}

	point, err := ed.GeFromBytes((*[32]byte)(&publicKey))
	if err != nil {
		return nil, err
	}

	var point3 ed.CompletedGroupElement
	point2 := ed.GeScalarMult((*[32]byte)(&secretKey), point)
	ed.GeMul8(&point3, &point2)
	point3.ToProjective(&point2)

	b := KeyDerivation(point2.ToBytes())
	return &b, nil
}

func GenerateKeyImage(publicKey PublicKey, secretKey SecretKey) (*KeyImage, error) {
	if !ed.ScCheck(secretKey) {
		return nil, errors.New("wrong secret key provided")
	}

	pkHash := HashFromBytes(publicKey[:])
	point, err := pkHash.ToEc()
	if err != nil {
		return nil, err
	}

	point2 := ed.GeScalarMult((*[32]byte)(&secretKey), point)
	keyImage := KeyImage(point2.ToBytes())

	return &keyImage, nil
}
