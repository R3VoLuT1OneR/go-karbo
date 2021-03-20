package cryptonote

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	TxTagCoinbase       = 0xff
	TxTagKey            = 0x2
	TxTagMultisignature = 0x3
)

type KeyImage EllipticCurvePointer

type InputCoinbase struct {
	Height uint32
}

type InputKey struct {
	Amount uint64
	OutputIndexes []uint32
	KeyImage
}

type InputMultisignature struct {
	Amount uint64
	SignatureCount uint8
	OutputIndex uint32
}

type OutputKey struct {
	Key
}

type OutputMultisignature struct {
	Keys []Key
	RequiredSignaturesCount byte
}

type TransactionInput interface {
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
	TransactionPrefix
}

type Transaction struct {
	TransactionPrefix

	Signature [][]Signature

	hash *Hash
	bytes *[]byte
}

func (t *Transaction) Serialize() ([]byte, error) {
	return t.TransactionPrefix.serialize()
}

func (t *Transaction) Deserialize(r *bytes.Reader) error {

	if err := t.TransactionPrefix.deserialize(r); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Hash() (*Hash, error) {
	if t.hash == nil {
		transactionBytes, err := t.Serialize()
		if err != nil {
			return nil, err
		}

		t.hash = new(Hash)
		t.hash.FromBytes(&transactionBytes)
	}

	return t.hash, nil
}

func (tp *TransactionPrefix) serialize() ([]byte, error) {
	var serialized bytes.Buffer

	varIntBuf := make([]byte, binary.MaxVarintLen64)

	written := binary.PutUvarint(varIntBuf, uint64(tp.Version))
	serialized.Write(varIntBuf[:written])

	written = binary.PutUvarint(varIntBuf, tp.UnlockHeight)
	serialized.Write(varIntBuf[:written])

	inputsLen := len(tp.Inputs)
	if inputsLen == 0 {
		return nil, errors.New("no inputs")
	}

	written = binary.PutUvarint(varIntBuf, uint64(inputsLen))
	serialized.Write(varIntBuf[:written])

	for _, input := range tp.Inputs {
		switch input.(type) {
		case InputCoinbase:
			serialized.WriteByte(TxTagCoinbase)

			written = binary.PutUvarint(varIntBuf, uint64(input.(InputCoinbase).Height))
			serialized.Write(varIntBuf[:written])
		default:
			return nil, errors.New(fmt.Sprintf("unknown input type: %T", input))
		}
	}

	outputLen := len(tp.Outputs)
	if outputLen == 0 {
		return nil, errors.New("no outputs")
	}

	written = binary.PutUvarint(varIntBuf, uint64(outputLen))
	serialized.Write(varIntBuf[:written])

	for _, output := range tp.Outputs {
		written = binary.PutUvarint(varIntBuf, output.Amount)
		serialized.Write(varIntBuf[:written])

		switch output.Target.(type) {
		case OutputKey:
			serialized.WriteByte(TxTagKey)
			serialized.Write(output.Target.(OutputKey).Bytes()[:])
		default:
			return nil, errors.New(fmt.Sprintf("unknown output target type: %T", output.Target))
		}
	}

	written = binary.PutUvarint(varIntBuf, uint64(len(tp.Extra)))
	serialized.Write(varIntBuf[:written])
	serialized.Write(tp.Extra[:])

	return serialized.Bytes(), nil
}

func (tp *TransactionPrefix) deserialize(r io.Reader) error {
	br := bufio.NewReader(r)

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
			// TODO: Implement transaction key
			return errors.New("not implemented")
		case TxTagMultisignature:
			// TODO: Implement multisig
			return errors.New("not implemented")
		default:
			return errors.New("unknown tx input tag")
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
			return errors.New("unknown tx output tag")
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
