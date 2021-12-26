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

// ScalarMult scalar multiplication of two images
func (image *KeyImage) ScalarMult(a *KeyImage) (*KeyImage, error) {
	A, err := ed.GeFromBytes((*[32]byte)(image))
	if err != nil {
		return nil, err
	}

	R := ed.GeScalarMult((*[32]byte)(a), A)

	aP := KeyImage(R.ToBytes())
	return &aP, nil
}
