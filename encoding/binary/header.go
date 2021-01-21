package binary

import "encoding/binary"

type storageBlockHeader struct {
	signatureA uint32
	signatureB uint32
	ver        byte
}

var baseHeadBlock = storageBlockHeader{
	signatureA: storageSignatureA,
	signatureB: storageSignatureB,
	ver: storageFormatVer,
}

func (h *storageBlockHeader) encode() [headSize]byte {
	var result [headSize]byte

	binary.LittleEndian.PutUint32(result[0:4], h.signatureA)
	binary.LittleEndian.PutUint32(result[4:8], h.signatureB)
	result[8] = h.ver

	return result
}

func (h *storageBlockHeader) decode(b [headSize]byte) error {
	h.signatureA = binary.LittleEndian.Uint32(b[0:4])
	h.signatureB = binary.LittleEndian.Uint32(b[4:8])
	h.ver = b[8]

	return nil
}

func (h *storageBlockHeader) equals(b storageBlockHeader) bool {
	return h.signatureA == b.signatureA && h.signatureB == b.signatureB && h.ver == b.ver
}