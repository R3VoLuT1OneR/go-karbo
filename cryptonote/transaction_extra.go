package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ErrPaddingNotZero = errors.New("padding byte is not zero")
	ErrPaddingMax = errors.New("padding size bigger than allowed")
	ErrNonceMax = errors.New("nonce size bigger than allowed")
	ErrMergeMiningTagMax = errors.New("merge mining tag is not 33")
)

var (
	TxExtraPaddingMax = 255
	TxExtraNonceMax = 255

	TxExtraTagPadding = byte(0x00)
	TxExtraTagPubkey = byte(0x01)
	TxExtraTagNonce = byte(0x02)
	TxExtraTagMergeMining = byte(0x03)

	TxExtraNoncePaymentId = 0x00

	TxExtranMergeMiningTagSize = uint64(33)
)

type TransactionExtraMergeMiningTag struct {
	Depth uint64
	MerkleRoot Hash
}

type TransactionExtraFields struct {
	PublicKey 	Key
	Nonce 		[]byte
	MiningTag   *TransactionExtraMergeMiningTag

	bytes []byte
}

func TxExtraFromBytes(b []byte) (*TransactionExtraFields, error) {
	r := bytes.NewReader(b)

	tef := TransactionExtraFields{}
	tef.bytes = b

	for {
		tag, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		switch tag {
		case TxExtraTagPadding:
			size := 0
			for {
				size++
				padding, err := r.ReadByte()
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}

				if padding != 0 {
					return nil, ErrPaddingNotZero
				}

				if size > TxExtraPaddingMax {
					return nil, ErrPaddingMax
				}
			}
		case TxExtraTagPubkey:
			kb := make([]byte, 32)
			if err := binary.Read(r, binary.LittleEndian, kb); err != nil {
				return nil, err
			}
			tef.PublicKey, err = KeyFromBytes(&kb)
			if err != nil {
				return nil, err
			}
		case TxExtraTagNonce:
			size, err := binary.ReadUvarint(r)
			if err != nil {
				return nil, err
			}
			if int(size) > TxExtraNonceMax {
				return nil, ErrNonceMax
			}

			if size > 0 {
				nb := make([]byte, size)
				if err := binary.Read(r, binary.LittleEndian, nb); err != nil {
					return nil, err
				}
				tef.Nonce = nb
			}
		case TxExtraTagMergeMining:
			size, err := binary.ReadUvarint(r)
			if err != nil {
				return nil, err
			}
			if size != TxExtranMergeMiningTagSize {
				return nil, ErrMergeMiningTagMax
			}

			depth, err := binary.ReadUvarint(r)
			if err != nil {
				return nil, err
			}

			var h Hash
			if err := h.Read(r); err != nil {
				return nil, err
			}

			tef.MiningTag = &TransactionExtraMergeMiningTag{depth, h}
		default:
			return nil, errors.New(fmt.Sprintf("Unknown extra tag: %x", tag))
		}
	}

	return &tef, nil
}

func (tp *TransactionPrefix) ParseExtra() (*TransactionExtraFields, error) {
	return TxExtraFromBytes(tp.Extra)
}
