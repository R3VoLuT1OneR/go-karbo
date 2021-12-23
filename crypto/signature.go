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
	h    Hash
	key  EllipticCurvePoint
	comm EllipticCurvePoint
}

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

	buf.h = hash
	buf.key = publicKey.Bytes()

tryAgain:
	k := RandomScalar()

	// we don't want tiny numbers here
	if binary.LittleEndian.Uint32(k[28:32]) == 0 {
		goto tryAgain
	}

	ed.GeScalarMultBase(&tmp3, k)
	buf.comm = tmp3.ToBytes()

	var bufBytes bytes.Buffer
	if err := binary.Write(&bufBytes, binary.LittleEndian, buf); err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	sig.C = HashToScalar(bufBytes.Bytes())
	if !ed.ScIsNonZero(sig.C) {
		goto tryAgain
	}

	sig.R = ed.ScMulSub(sig.C, secretKey.Bytes(), k)
	if !ed.ScIsNonZero(sig.R) {
		goto tryAgain
	}

	return &sig, nil
}
func (s *Signature) Check(h Hash, pk Key) bool {

	return false
}

func (s *Signature) Deserialize(r *bytes.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, s); err != nil {
		return err
	}

	return nil
}

func (s *Signature) Serialize() ([]byte, error) {
	var serialized bytes.Buffer

	if err := binary.Write(&serialized, binary.LittleEndian, s); err != nil {
		return nil, err
	}

	return serialized.Bytes(), nil
}
