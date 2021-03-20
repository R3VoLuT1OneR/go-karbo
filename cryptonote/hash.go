package cryptonote

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
)

type Hash [32]byte

type HashList []Hash

func (h *Hash) FromBytes(b *[]byte) {
	hashed := hash.Keccak(b)
	copy(h[:32], hashed[:32])
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
	hashed := hash.Keccak(&b)
	var h Hash
	copy(h[:32], hashed[:32])
	return h
}

func (hl HashList) merkleRootHash() (*Hash, error)  {
	switch len(hl) {
	case 0:
		return nil, errors.New("at least 1 hash must be provided")
	case 1:
		singleHash := hl[0]
		return &singleHash, nil
	case 2:
		h := HashFromBytes(append(hl[0][:], hl[1][:]...))
		return &h, nil
	default:
		cnt := 2 // Largest power of two
		for cnt << 1 < len(hl) {
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

		return &tempList[0], nil
	}
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
