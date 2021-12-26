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

type ringBuf struct {
	hash *Hash
	sigs []Signature
}

func (buf *ringBuf) toScalar() (*EllipticCurveScalar, error) {
	var bytesBuf bytes.Buffer

	bytesBuf.Write(buf.hash[:])

	for i, signature := range buf.sigs {
		sigBytes, err := signature.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize signature at %d: %w", i, err)
		}

		bytesBuf.Write(sigBytes)
	}

	bufHash := HashFromBytes(bytesBuf.Bytes())
	scalar := bufHash.toScalar()

	return &scalar, nil
}

// GenerateRingSignature
// Procedure generate_signature(M, A[1], A[2], ..., A[n], i, a[i]):
//   I <- a[i]*H(A[i])
//   c[j], r[j] [j=1..n, j!=i] <- random
//   k <- random
//   For j <- 1..n, j!=i
//     X[j] <- c[j]*A[j]+r[j]*G
//     Y[j] <- c[j]*I+r[j]*H(A[j])
//   End For
//   X[i] <- k*G
//   Y[i] <- k*H(A[i])
//   c[i] <- H(H(M) || X[1] || Y[1] || X[2] || Y[2] || ... || X[n] ||
//   Y[n])-Sum[j=1..n, j!=i](c[j])
//   r[i] <- k-a[i]*c[i]
// Return (I, c[1] || r[1] || c[2] || r[2] || ... || c[n] || r[n])
// End Procedure

func GenerateRingSignature(prefixHash *Hash, image *KeyImage, pubs *[]PublicKey, sec *SecretKey, secIndex uint64) ([]Signature, error) {
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

	if *image != *t3 {
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
	sum := ed.ScZero()

	buf := ringBuf{
		hash: prefixHash,
		sigs: make([]Signature, len(*pubs)),
	}

	k := RandomScalar()
	sigs := make([]Signature, len(*pubs))
	for i := 0; i < len(buf.sigs); i++ {
		var tmp2 ed.ProjectiveGroupElement
		var tmp3 ed.ExtendedGroupElement

		if uint64(i) == secIndex {
			ed.GeScalarMultBase(&tmp3, (*[32]byte)(&k))
			buf.sigs[i].C = tmp3.ToBytes()

			pubHash := HashFromBytes(((*pubs)[i])[:])
			tmp3, err := pubHash.toEc()
			if err != nil {
				return nil, fmt.Errorf("failed to get EC from pub hash at %d: %w", i, err)
			}
			tmp2 = ed.GeScalarMult((*[32]byte)(&k), tmp3)
			buf.sigs[i].R = tmp2.ToBytes()
			continue
		}

		sigs[i].C = RandomScalar()
		sigs[i].R = RandomScalar()

		tmp3p, err := ed.GeFromBytes((*[32]byte)(&(*pubs)[i]))
		if err != nil {
			return nil, fmt.Errorf("failed to get tmp3 from bytes at %d: %w", i, err)
		}
		tmp3 = *tmp3p

		ed.GeDoubleScalarMultBaseVartime(&tmp2, (*[32]byte)(&sigs[i].C), &tmp3, (*[32]byte)(&sigs[i].R))
		buf.sigs[i].C = tmp2.ToBytes()

		pubHash := HashFromBytes(((*pubs)[i])[:])
		tmp3p, err = pubHash.toEc()
		if err != nil {
			return nil, fmt.Errorf("failed to get EC from pub hash at %d: %w", i, err)
		}
		tmp3 = *tmp3p

		ed.GeDoubleScalarMultPrecompVartime(&tmp2, (*[32]byte)(&sigs[i].R), &tmp3, (*[32]byte)(&sigs[i].C), imagePre)
		buf.sigs[i].R = tmp2.ToBytes()

		sum = ed.ScAdd(&sum, (*[32]byte)(&sigs[i].C))
	}

	h, err := buf.toScalar()
	if err != nil {
		return nil, err
	}

	sigs[secIndex].C = ed.ScSub(*h, sum)
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

	checkResult := ed.GeCheckSubGroupPreCompVartime(imagePre)
	if checkKeyImage && checkResult != 0 {
		return false
	}

	buf := ringBuf{
		hash: prefixHash,
		sigs: make([]Signature, len(*pubs)),
	}

	sum := ed.ScZero()
	for i := 0; i < len(*pubs); i++ {
		var tmp2 ed.ProjectiveGroupElement
		var tmp3 ed.ExtendedGroupElement
		var sig = &((*sigs)[i])
		var pub = &((*pubs)[i])

		if !ed.ScCheck(sig.C) || !ed.ScCheck(sig.R) {
			return false
		}

		if !tmp3.FromBytes((*[32]byte)(pub)) {
			return false
		}

		ed.GeDoubleScalarMultBaseVartime(&tmp2, (*[32]byte)(&sig.C), &tmp3, (*[32]byte)(&sig.R))
		buf.sigs[i].C = tmp2.ToBytes()

		pubHash := HashFromBytes(pub[:])
		tmp3p, err := pubHash.toEc()
		if err != nil {
			return false
		}
		tmp3 = *tmp3p

		ed.GeDoubleScalarMultPrecompVartime(&tmp2, (*[32]byte)(&sig.R), &tmp3, (*[32]byte)(&sig.C), imagePre)
		buf.sigs[i].R = tmp2.ToBytes()

		sum = ed.ScAdd(&sum, (*[32]byte)(&sig.C))
	}

	h, err := buf.toScalar()
	if err != nil {
		return false
	}

	return !ed.ScIsNonZero(ed.ScSub(*h, sum))
}

func (e EllipticCurveScalar) bytesPointer() *[32]byte {
	var b [32]byte = e
	return &b
}
