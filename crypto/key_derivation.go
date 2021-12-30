package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

type KeyDerivation EllipticCurvePoint

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

func (derivation *KeyDerivation) toPublicKey(outputIndex uint64, base *PublicKey) (*PublicKey, error) {
	point1, err := ed.GeFromBytes((*[32]byte)(base))
	if err != nil {
		return nil, err
	}

	scalar := derivation.toScalar(outputIndex)

	var point2 ed.ExtendedGroupElement
	ed.GeScalarMultBase(&point2, (*[32]byte)(&scalar))

	var point3 ed.CachedGroupElement
	point2.ToCached(&point3)

	var point4 ed.CompletedGroupElement
	ed.GeAdd(&point4, point1, &point3)

	var point5 ed.ProjectiveGroupElement
	point4.ToProjective(&point5)

	publicKey := PublicKey(point5.ToBytes())
	return &publicKey, err
}

func (derivation *KeyDerivation) toSecretKey(outputIndex uint64, base *PublicKey) (*SecretKey, error) {
	if !ed.ScCheck(*base) {
		return nil, errors.New("base is wrong public key")
	}

	scalar := derivation.toScalar(outputIndex)
	secretKey := SecretKey(ed.ScAdd((*[32]byte)(base), (*[32]byte)(&scalar)))

	return &secretKey, nil
}

func (derivation *KeyDerivation) toScalar(outputIndex uint64) EllipticCurveScalar {
	var b bytes.Buffer

	outputIndexBytes := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(outputIndexBytes, outputIndex)

	b.Write(derivation[:])
	b.Write(outputIndexBytes[:written])

	hash := HashFromBytes(b.Bytes())

	return hash.toScalar()
}
