package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
)

type Hash [32]byte

type HashList []Hash

func (h *Hash) FromBytes(b *[]byte) {
	hashed := hash.Keccak(b)
	copy(h[:32], hashed[:32])
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
		//hashedHashList := HashList{}
		//
		//for i := 0; i < (listLen - (listLen % 2)); i += 2 {
		//	thl := HashList{hl[i], hl[i+1]}
		//
		//	if listLen - i == 3 {
		//		lastHL := HashList{hl[listLen-2], hl[listLen-1]}
		//		lastH, err := lastHL.merkleRootHash()
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		thl[1] = lastH
		//	}
		//
		//	th, err := thl.merkleRootHash()
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	hashedHashList = append(hashedHashList, th)
		//}

		//hashedHashListLen := listLen / 2
		//hashedHashList := HashList{}
		//
		//for i := 0; i < hashedHashListLen; i++ {
		//	thl := HashList{hl[i*2], hl[i*2+1]}
		//
		//	th, err := thl.merkleRootHash()
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	hashedHashList = append(hashedHashList, th)
		//}
		//
		//if listLen > hashedHashListLen * 2 {
		//	thl := HashList{hl[listLen - 2], hl[listLen-1]}
		//	lh, err := thl.merkleRootHash()
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	thl = HashList{hl[listLen - 3], lh}
		//	th, err := thl.merkleRootHash()
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	hashedHashList[len(hashedHashList) - 1] = th
		//}
		//
		//return hashedHashList.merkleRootHash()
	}
}
