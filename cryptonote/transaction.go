package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	TxTagCoinbase       byte = 0xff
	TxTagKey                 = 0x2
	TxTagMultisignature      = 0x3
)

type KeyImage EllipticCurvePointer

type TransactionSignatures [][]Signature

type InputCoinbase struct {
	BlockIndex uint32
}

func (i InputCoinbase) sigCount() int {
	return 0
}

type InputKey struct {
	Amount        uint64
	OutputIndexes []uint32
	KeyImage
}

func (i InputKey) sigCount() int {
	return len(i.OutputIndexes)
}

type InputMultisignature struct {
	Amount         uint64
	SignatureCount uint8
	OutputIndex    uint32
}

func (i InputMultisignature) sigCount() int {
	return int(i.SignatureCount)
}

type OutputKey struct {
	Key
}

type OutputMultisignature struct {
	Keys                    []Key
	RequiredSignaturesCount byte
}

type TransactionInput interface {
	sigCount() int
}

type TransactionOutputTarget interface {
}

type TransactionOutput struct {
	Amount uint64
	Target TransactionOutputTarget
}

type TransactionPrefix struct {
	Version      byte
	UnlockHeight uint64
	Inputs       []TransactionInput
	Outputs      []TransactionOutput
	Extra        []byte
}

type BaseTransaction struct {
	*TransactionPrefix

	hash *Hash
}

type Transaction struct {
	TransactionPrefix
	TransactionSignatures

	hash *Hash
}

func (t *Transaction) Serialize() []byte {
	pb := t.TransactionPrefix.serialize()
	sb := t.TransactionSignatures.serialize(t)

	return append(pb, sb...)
}

func (t *Transaction) Deserialize(r *bytes.Reader) error {
	if err := t.TransactionPrefix.deserialize(r); err != nil {
		return err
	}

	if err := t.TransactionSignatures.deserialize(r, t); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Size() uint64 {
	return uint64(len(t.Serialize()))
}

func (t *Transaction) Hash() *Hash {
	if t.hash == nil {
		t.hash = new(Hash)
		t.hash.FromBytes(t.Serialize())
	}

	return t.hash
}

func (t *BaseTransaction) Hash() *Hash {
	if t.hash == nil {
		t.hash = new(Hash)
		t.hash.FromBytes(t.TransactionPrefix.serialize())
	}

	return t.hash
}

func (tp *TransactionPrefix) serialize() []byte {
	var serialized bytes.Buffer

	varIntBuf := make([]byte, binary.MaxVarintLen64)

	written := binary.PutUvarint(varIntBuf, uint64(tp.Version))
	serialized.Write(varIntBuf[:written])

	written = binary.PutUvarint(varIntBuf, tp.UnlockHeight)
	serialized.Write(varIntBuf[:written])

	written = binary.PutUvarint(varIntBuf, uint64(len(tp.Inputs)))
	serialized.Write(varIntBuf[:written])

	for _, input := range tp.Inputs {
		switch input.(type) {
		case InputCoinbase:
			serialized.WriteByte(TxTagCoinbase)

			written = binary.PutUvarint(varIntBuf, uint64(input.(InputCoinbase).BlockIndex))
			serialized.Write(varIntBuf[:written])
		case InputKey:
			inputKey := input.(InputKey)
			serialized.WriteByte(TxTagKey)

			// Write amount
			written = binary.PutUvarint(varIntBuf, inputKey.Amount)
			serialized.Write(varIntBuf[:written])

			// Write output indexes
			size := len(inputKey.OutputIndexes)
			written = binary.PutUvarint(varIntBuf, uint64(size))
			serialized.Write(varIntBuf[:written])

			for _, outputIndex := range inputKey.OutputIndexes {
				written = binary.PutUvarint(varIntBuf, uint64(outputIndex))
				serialized.Write(varIntBuf[:written])
			}

			// Write key image
			_ = binary.Write(&serialized, binary.LittleEndian, inputKey.KeyImage)
		default:
		}
	}

	written = binary.PutUvarint(varIntBuf, uint64(len(tp.Outputs)))
	serialized.Write(varIntBuf[:written])

	for _, output := range tp.Outputs {
		written = binary.PutUvarint(varIntBuf, output.Amount)
		serialized.Write(varIntBuf[:written])

		switch output.Target.(type) {
		case OutputKey:
			serialized.WriteByte(TxTagKey)
			serialized.Write(output.Target.(OutputKey).Bytes()[:])
		default:
		}
	}

	written = binary.PutUvarint(varIntBuf, uint64(len(tp.Extra)))
	serialized.Write(varIntBuf[:written])
	serialized.Write(tp.Extra[:])

	return serialized.Bytes()
}

func (tp *TransactionPrefix) deserialize(br *bytes.Reader) error {
	/**
	 * Read transaction version
	 */
	version, err := binary.ReadUvarint(br)
	if err != nil {
		return err
	}
	tp.Version = byte(version)

	/**
	 * Read transaction UnlockHeight
	 */
	tp.UnlockHeight, err = binary.ReadUvarint(br)
	if err != nil {
		return err
	}

	/**
	 * Read transaction Inputs
	 */
	inputsLen, err := binary.ReadUvarint(br)
	tp.Inputs = make([]TransactionInput, inputsLen)
	if err != nil {
		return err
	}

	for inputIndex := uint64(0); inputIndex < inputsLen; inputIndex++ {
		var tag byte
		if err := binary.Read(br, binary.LittleEndian, &tag); err != nil {
			return err
		}

		switch tag {
		case TxTagCoinbase:
			blockIndex, err := binary.ReadUvarint(br)
			if err != nil {
				return err
			}

			tp.Inputs[inputIndex] = InputCoinbase{uint32(blockIndex)}
		case TxTagKey:
			var OutputIndexes []uint32
			var Key KeyImage

			amount, err := binary.ReadUvarint(br)
			if err != nil {
				return err
			}

			size, err := binary.ReadUvarint(br)
			if err != nil {
				return err
			}

			for i := uint64(0); i < size; i++ {
				oi, err := binary.ReadUvarint(br)
				if err != nil {
					return err
				}

				OutputIndexes = append(OutputIndexes, uint32(oi))
			}

			if err := binary.Read(br, binary.LittleEndian, &Key); err != nil {
				return err
			}

			tp.Inputs[inputIndex] = InputKey{
				Amount:        amount,
				OutputIndexes: OutputIndexes,
				KeyImage:      Key,
			}
		case TxTagMultisignature:
			// TODO: Implement multisig
			return errors.New("not implemented")
		default:
			return errors.New(fmt.Sprintf("unknown tx input tag: %x", tag))
		}
	}

	/**
	 * Read transaction Output
	 */
	outputLen, err := binary.ReadUvarint(br)
	tp.Outputs = make([]TransactionOutput, outputLen)
	if err != nil {
		return err
	}

	for outputIndex := uint64(0); outputIndex < outputLen; outputIndex++ {
		amount, err := binary.ReadUvarint(br)
		if err != nil {
			return err
		}

		tag, err := binary.ReadUvarint(br)
		if err != nil {
			return err
		}

		switch tag {
		case TxTagKey:
			var keyBytes [32]byte
			if err := binary.Read(br, binary.LittleEndian, &keyBytes); err != nil {
				return err
			}

			tp.Outputs[outputIndex] = TransactionOutput{
				Amount: amount,
				Target: OutputKey{KeyFromArray(&keyBytes)},
			}
		case TxTagMultisignature:
			// TODO: Implement multisig
			return errors.New("not implemented")
		default:
			return errors.New(fmt.Sprintf("unknown tx output tag: %x", tag))
		}
	}

	/**
	 * Read transaction Extra
	 */
	extraLen, err := binary.ReadUvarint(br)
	if err != nil {
		return err
	}

	tp.Extra = make([]byte, extraLen)
	if err := binary.Read(br, binary.LittleEndian, &tp.Extra); err != nil {
		return err
	}

	return nil
}

func (ts TransactionSignatures) serialize(t *Transaction) []byte {
	var serialized bytes.Buffer

	for i, input := range t.TransactionPrefix.Inputs {
		sigSize := input.sigCount()

		if sigSize == 0 {
			continue
		}

		for _, sig := range ts[i] {
			_ = binary.Write(&serialized, binary.LittleEndian, sig)
		}
	}

	return serialized.Bytes()
}

func (ts *TransactionSignatures) deserialize(br *bytes.Reader, t *Transaction) error {
	inputs := t.TransactionPrefix.Inputs
	signaturesNotExpected := len(inputs) == 0

	if len(inputs) == 1 {
		if _, ok := inputs[0].(InputCoinbase); ok {
			signaturesNotExpected = true
		}
	}

	for _, input := range inputs {
		sigSize := input.sigCount()

		if signaturesNotExpected && sigSize != 0 {
			return errors.New("unexpected signatures")
		}

		if sigSize == 0 {
			continue
		}

		sigs := make([]Signature, sigSize)
		for i := 0; i < sigSize; i++ {
			if err := binary.Read(br, binary.LittleEndian, &sigs[i]); err != nil {
				return err
			}
		}

		*ts = append(*ts, sigs)
	}

	return nil
}
