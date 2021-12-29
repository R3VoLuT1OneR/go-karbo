package crypto

import (
	"crypto/rand"
	"errors"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
	"io"
)

type EllipticCurvePoint [32]byte

type EllipticCurveScalar [32]byte

type KeyDerivation EllipticCurvePoint

type KeyImage EllipticCurvePoint

// Check that the point is on curve
func (p *EllipticCurvePoint) Check() bool {
	var point ed.ExtendedGroupElement
	return point.FromBytes((*[32]byte)(p))
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

//func DerivePublicKey(derivation *KeyDerivation, outputIndex int, base *PublicKey) (*PublicKey, error) {
//	point1, err := ed.GeFromBytes((*[32]byte)(base))
//	if err != nil {
//		return nil, err
//	}
//
//}
//
//func (derivation *KeyDerivation) toScalar(outputIndex int) (EllipticCurveScalar, error) {
//	//var sizeOfInt = (int)(*(*uint)(unsafe.Sizeof(outputIndex)))
//	type bufType struct {
//		derivation  KeyDerivation
//		outputIndex []byte
//	}
//
//	buf := bufType{
//		derivation: *derivation,
//	}
//
//}
