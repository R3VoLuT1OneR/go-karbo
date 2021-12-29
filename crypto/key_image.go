package crypto

import (
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

type KeyImage EllipticCurvePoint

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
