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

	ed.GeScalarMultBase(&tmp3, (*[32]byte)(&k))
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
	ed.GeDoubleScalarMultBaseVartime(&tmp2, sig.C.bytesPointer(), &tmp3, sig.R.bytesPointer())
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

func GenerateRingSignature(prefixHash Hash, image *KeyImage, pubs *[]PublicKey, sec *SecretKey, secIndex uint64) ([]Signature, error) {
	var sigs []Signature

	if !(secIndex < uint64(len(*pubs))) {
		return nil, errors.New("assertion failed: wrong secIndex provided")
	}

	// -------- ENABLED IN DEBUG MODE --------
	if !ed.ScCheck(*sec) {
		return nil, errors.New("wrong secret key provided")
	}

	var t ed.ExtendedGroupElement
	ed.GeScalarMultBase(&t, (*[32]byte)(sec))

	t2 := PublicKey(t.ToBytes())

	secIndexPub := (*pubs)[secIndex]
	if secIndexPub != t2 {
		return nil, errors.New("assertion failed: sec index pub key not match")
	}

	t3, err := GenerateKeyImage(&secIndexPub, sec)
	if err != nil {
		return nil, err
	}

	if image != t3 {
		return nil, errors.New("assertion failed: wrong image provided")
	}

	for i := 0; i < len(*pubs); i++ {
		if !(*pubs)[i].Check() {
			return nil, errors.New(fmt.Sprintf("wrong public key at index %d", i))
		}
	}
	// -------- ENABLED IN DEBUG MODE --------

	imageUnp, err := ed.GeFromBytes((*[32]byte)(image))
	if err != nil {
		return nil, err
	}

	imagePre := ed.GeDSMPreComp(imageUnp)
	sum := [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	k := RandomScalar()
	bufSigs := make([]Signature, len(*pubs))
	for i := 0; i < len(bufSigs); i++ {
		var tmp2 *ed.ProjectiveGroupElement
		var tmp3 *ed.ExtendedGroupElement

		if uint64(i) == secIndex {
			ed.GeScalarMultBase(tmp3, (*[32]byte)(&k))
			bufSigs[i].C = tmp3.ToBytes()

			pubHash := HashFromBytes(((*pubs)[i])[:])
			tmp3, err := pubHash.toEc()
			if err != nil {
				return nil, fmt.Errorf("failed to get EC from pub hash at %d: %w", i, err)
			}
			*tmp2 = ed.GeScalarMult((*[32]byte)(&k), tmp3)
			bufSigs[i].R = tmp2.ToBytes()
			continue
		}

		sigs[i].C = RandomScalar()
		sigs[i].R = RandomScalar()

		tmp3, err := ed.GeFromBytes((*[32]byte)(&(*pubs)[i]))
		if err != nil {
			return nil, fmt.Errorf("failed to get tmp3 from bytes at %d: %w", i, err)
		}

		ed.GeDoubleScalarMultBaseVartime(tmp2, (*[32]byte)(&sigs[i].C), tmp3, (*[32]byte)(&sigs[i].R))
		bufSigs[i].C = tmp2.ToBytes()

		pubHash := HashFromBytes(((*pubs)[i])[:])
		tmp3, err = pubHash.toEc()
		if err != nil {
			return nil, fmt.Errorf("failed to get EC from pub hash at %d: %w", i, err)
		}

		ed.GeDoubleScalarMultPrecompVartime(tmp2, (*[32]byte)(&sigs[i].R), tmp3, (*[32]byte)(&sigs[i].C), imagePre)
		bufSigs[i].R = tmp2.ToBytes()

		sum = ed.ScAdd(&sum, (*[32]byte)(&sigs[i].C))
	}

	// Serialize prefix hash and new generated signatures
	var buf bytes.Buffer
	buf.Write(prefixHash[:])
	for i, signature := range bufSigs {
		sigBytes, err := signature.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize signature at %d: %w", i, err)
		}

		buf.Write(sigBytes)
	}

	bufHash := HashFromBytes(buf.Bytes())
	h := bufHash.toScalar()

	sigs[secIndex].C = ed.ScSub(h, sum)
	sigs[secIndex].R = ed.ScMulSub(sigs[secIndex].C, *sec, k)

	return sigs, nil
}

func CheckRingSignature(prefixHash *Hash, image *KeyImage, pubs *[]PublicKey, sigs *[]Signature, checkKeyImage bool) bool {
	// -------- ENABLED IN DEBUG MODE --------
	for i := 0; i < len(*pubs); i++ {
		if !(*pubs)[i].Check() {
			return false
		}
	}
	// -------- ENABLED IN DEBUG MODE --------

	imageUnp, err := ed.GeFromBytes((*[32]byte)(image))
	if err != nil {
		return false
	}

	imagePre := ed.GeDSMPreComp(imageUnp)

}

func (e EllipticCurveScalar) bytesPointer() *[32]byte {
	var b [32]byte = e
	return &b
}
