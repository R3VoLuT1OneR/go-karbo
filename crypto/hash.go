package crypto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

const HashLength = 32

type Hash [HashLength]byte

type HashList []Hash

func (h *Hash) FromBytes(b []byte) {
	hashed := Keccak(b)
	copy(h[:HashLength], hashed[:HashLength])
}

func (h *Hash) Read(r *bytes.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, h); err != nil {
		return err
	}

	return nil
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
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
