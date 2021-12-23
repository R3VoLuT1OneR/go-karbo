package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	if err := binary.Write(&bufBytes, binary.LittleEndian, *s); err != nil {
		return nil, err
	}

	return bufBytes.Bytes(), nil
}

// GenerateSignature generates new signature
func GenerateSignature(hash Hash, publicKey Key, secretKey Key) (*Signature, error) {
	sig := Signature{}

	var tmp3 ed.ExtendedGroupElement
	var buf sComm

	// Check that provided public key belongs to secret key
	tPubKey, err := PublicFromPrivate(secretKey)
	if err != nil {
		return nil, fmt.Errorf("can't get public from private: %w", err)
	}

	if tPubKey.Bytes() != publicKey.Bytes() {
		return nil, errors.New("mismatch in provided public and secret keys")
	}

	buf.hash = hash
	buf.key = publicKey.Bytes()

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

	sig.C = HashToScalar(bufBytes)
	if !ed.ScIsNonZero(sig.C) {
		goto tryAgain
	}

	sig.R = ed.ScMulSub(sig.C, secretKey.Bytes(), k)
	if !ed.ScIsNonZero(sig.R) {
		goto tryAgain
	}

	return &sig, nil
}

func (sig *Signature) Check(hash Hash, publicKey Key) bool {
	var tmp2 ed.ProjectiveGroupElement
	var tmp3 ed.ExtendedGroupElement
	var buf sComm

	if !publicKey.Check() {
		return false
	}

	publicKeyBytes := publicKey.Bytes()

	buf.hash = hash
	buf.key = publicKeyBytes

	if !tmp3.FromBytes(&publicKeyBytes) {
		return false
	}

	if !ed.ScCheck(sig.C) || !ed.ScCheck(sig.R) || !ed.ScIsNonZero(sig.C) {
		return false
	}

	ed.GeDoubleScalarMultVartime(&tmp2, sig.C.bytesPointer(), &tmp3, sig.R.bytesPointer())
	buf.comm = tmp2.ToBytes()

	if bytes.Compare(buf.comm[:], infinityPoint[:]) == 0 {
		return false
	}

	bufBytes, err := buf.bytes()
	if err != nil {
		return false
	}

	c := HashToScalar(bufBytes)
	c = ed.ScSub(c, sig.C)

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
