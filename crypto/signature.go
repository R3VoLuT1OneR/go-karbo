package crypto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

type Signature struct {
	C, R EllipticCurveScalar
}

type sComm struct {
	hash Hash
	key  EllipticCurvePoint
	comm EllipticCurvePoint
}

var infinityPoint = EllipticCurvePoint{
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func (s *sComm) bytes() ([]byte, error) {
	var bufBytes bytes.Buffer
	if err := binary.Write(&bufBytes, binary.BigEndian, *s); err != nil {
		return nil, err
	}

	return bufBytes.Bytes(), nil
}

// Sign the hash with the private key.
func (hash Hash) Sign(secretKey *SecretKey) (*Signature, error) {
	sig := Signature{}

	var tmp3 ed.ExtendedGroupElement
	var buf sComm

	publicKey, err := PublicFromSecret(secretKey)
	if err != nil {
		return nil, fmt.Errorf("can't get public from private: %w", err)
	}

	buf.hash = hash
	buf.key = EllipticCurvePoint(*publicKey)

tryAgain:
	k := RandomScalar()

	// we don't want tiny numbers here
	if binary.LittleEndian.Uint32(k[28:32]) == 0 {
		goto tryAgain
	}

	ed.GeScalarMultBase(&tmp3, k)
	buf.comm = tmp3.ToBytes()
	bufBytes, err := buf.bytes()
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	bufHash := HashFromBytes(bufBytes)
	sig.C = bufHash.toScalar()
	if !ed.ScIsNonZero(sig.C) {
		goto tryAgain
	}

	sig.R = ed.ScMulSub(sig.C, *secretKey, k)
	if !ed.ScIsNonZero(sig.R) {
		goto tryAgain
	}

	return &sig, nil
}

func (sig *Signature) Check(hash *Hash, publicKey *PublicKey) bool {
	if !publicKey.Check() {
		return false
	}

	buf := sComm{
		hash: *hash,
		key:  EllipticCurvePoint(*publicKey),
	}

	var tmp3 ed.ExtendedGroupElement
	if !tmp3.FromBytes((*[32]byte)(publicKey)) {
		return false
	}

	if !ed.ScCheck(sig.C) || !ed.ScCheck(sig.R) || !ed.ScIsNonZero(sig.C) {
		return false
	}

	var tmp2 ed.ProjectiveGroupElement
	ed.GeDoubleScalarMultVartime(&tmp2, sig.C.bytesPointer(), &tmp3, sig.R.bytesPointer())
	buf.comm = tmp2.ToBytes()

	if bytes.Compare(buf.comm[:], infinityPoint[:]) == 0 {
		return false
	}

	bufBytes, _ := buf.bytes()
	bufHash := HashFromBytes(bufBytes)

	c := ed.ScSub(sig.C, bufHash.toScalar())
	return !ed.ScIsNonZero(c)
}

func (sig *Signature) Deserialize(r *bytes.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, sig); err != nil {
		return err
	}

	return nil
}

func (sig *Signature) Serialize() ([]byte, error) {
	var serialized bytes.Buffer

	if err := binary.Write(&serialized, binary.LittleEndian, sig); err != nil {
		return nil, err
	}

	return serialized.Bytes(), nil
}

func (e EllipticCurveScalar) bytesPointer() *[32]byte {
	var b [32]byte = e
	return &b
}
