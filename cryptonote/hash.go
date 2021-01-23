package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
)

type Hash [32]byte

type hashList []*Hash

func (h *Hash) FromBytes(b *[]byte) {
	hashed := hash.Keccak(b)
	copy(h[:32], hashed[:32])
}

func HashFromBytes(b *[]byte) Hash {
	hashed := hash.Keccak(b)
	var h Hash
	copy(h[:32], hashed[:32])
	return h
}

// TODO: Add tests from C++ implementation. tests/Hash/tests-tree.txt
func (hl hashList) merkleRootHash() (*Hash, error)  {
	var h Hash
	listLen := len(hl)

	switch listLen {
	case 0:
		return nil, errors.New("at least 1 hash must be provided")
	case 1:
		singleHash := hl[0]
		return singleHash, nil
	case 2:
		doubleHash := hl[0][:]
		doubleHash = append(doubleHash, hl[1][:]...)
		h.FromBytes(&doubleHash)
	default:
		hashedHashListLen := listLen / 2
		hashedHashList := hashList{}

		for i := 0; i < hashedHashListLen; i++ {
			thl := hashList{hl[i*2], hl[i*2+1]}

			th, err := thl.merkleRootHash()
			if err != nil {
				return nil, err
			}

			hashedHashList = append(hashedHashList, th)
		}

		if listLen > hashedHashListLen * 2 {
			thl := hashList{hl[listLen - 2], hl[listLen-1]}
			lh, err := thl.merkleRootHash()
			if err != nil {
				return nil, err
			}

			thl = hashList{hl[listLen - 3], lh}
			th, err := thl.merkleRootHash()
			if err != nil {
				return nil, err
			}

			hashedHashList[len(hashedHashList) - 1] = th
		}


		return hashedHashList.merkleRootHash()
	}

	return &h, nil
}
