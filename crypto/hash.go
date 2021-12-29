package crypto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	ed "github.com/r3volut1oner/go-karbo/crypto/edwards25519"
)

const HashLength = 32

type Hash [HashLength]byte

type HashList []Hash

func (hash *Hash) FromBytes(b []byte) {
	hashed := Keccak(b)
	copy(hash[:HashLength], hashed[:HashLength])
}

func (hash *Hash) Read(r *bytes.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, hash); err != nil {
		return err
	}

	return nil
}

func (hash *Hash) String() string {
	return hex.EncodeToString(hash[:])
}

func (hash *Hash) toEc() (*ed.ExtendedGroupElement, error) {
	var p ed.ProjectiveGroupElement
	var p2 ed.CompletedGroupElement
	var r ed.ExtendedGroupElement

	if err := p.FromBytes((*[32]byte)(hash)); err != nil {
		return nil, err
	}

	ed.GeMul8(&p2, &p)
	p2.ToExtended(&r)

	return &r, nil
}

func (hash *Hash) toPoint() (*EllipticCurvePoint, error) {
	r, err := hash.toEc()
	if err != nil {
		return nil, err
	}

	ecPoint := EllipticCurvePoint(r.ToBytes())

	return &ecPoint, nil
}

func (hash *Hash) toScalar() EllipticCurveScalar {
	return ed.ScReduce32(*hash)
}

func HashFromBytes(b []byte) Hash {
	hashed := Keccak(b)
	var h Hash
	copy(h[:32], hashed[:32])
	return h
}

func (hl HashList) MerkleRootHash() *Hash {
	switch len(hl) {
	case 0:
		// return nil, errors.New("at least 1 hash must be provided")
	case 1:
		singleHash := hl[0]
		return &singleHash
	case 2:
		h := HashFromBytes(append(hl[0][:], hl[1][:]...))
		return &h
	default:
		cnt := 2 // Largest power of two
		for cnt<<1 < len(hl) {
			cnt <<= 1
		}

		readyNum := (2 * cnt) - len(hl)
		tempList := make(HashList, readyNum)
		copy(tempList, hl[:readyNum])

		for i, j := readyNum, readyNum; j < cnt; i, j = i+2, j+1 {
			h := HashFromBytes(append(hl[i][:], hl[i+1][:]...))
			tempList = append(tempList, h)
		}

		for len(tempList) > 1 {
			newTempList := HashList{}
			for i := 0; i < len(tempList); i += 2 {
				h := HashFromBytes(append(tempList[i][:], tempList[i+1][:]...))
				newTempList = append(newTempList, h)
			}
			tempList = newTempList
		}

		return &tempList[0]
	}

	return nil
}

func (hl HashList) TreeHashFromBranch(leaf Hash) Hash {
	depth := len(hl)

	if depth == 0 {
		return leaf
	}

	fromLeaf := true

	var buf [2][32]byte
	var leafPath, branchPath *[32]byte

	for depth > 0 {
		depth--

		// TODO: WTF?
		//if (path && (((const char *) path)[depth >> 3] & (1 << (depth & 7))) != 0) {
		//	leaf_path = buffer[1];
		//	branch_path = buffer[0];
		//} else {
		//	leaf_path = buffer[0];
		//	branch_path = buffer[1];
		//}

		leafPath = &buf[0]
		branchPath = &buf[1]

		if fromLeaf {
			copy(leafPath[:], leaf[:])
			fromLeaf = false
		} else {
			h := HashFromBytes(append(buf[0][:], buf[1][:]...))
			copy(leafPath[:], h[:])
		}

		copy(branchPath[:], hl[depth][:])
	}

	return HashFromBytes(append(buf[0][:], buf[1][:]...))
}

func (hl *HashList) Index(h *Hash) int {
	for i, th := range *hl {
		if th == *h {
			return i
		}
	}

	return -1
}

func (hl *HashList) Has(h *Hash) bool {
	return hl.Index(h) >= 0
}

func (hl *HashList) Remove(h *Hash) {
	i := hl.Index(h)
	if i >= 0 {
		l := *hl
		l = append(l[:i], l[i+1:]...)
		*hl = l
	}
}
