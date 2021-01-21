package binary

import (
	"bytes"
	"errors"
)

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	headBlockBytes := baseHeadBlock.encode()
	buf.Write(headBlockBytes[:])

	encoder := NewEncoder(&buf)
	if err := encoder.Encode(v, ""); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Unmarshal(b []byte, v interface{}) error {
	reader := bytes.NewReader(b)

	var headBytes [headSize]byte
	var headBlock storageBlockHeader
	if _, err := reader.Read(headBytes[:]); err != nil {
		return err
	}

	if err := headBlock.decode(headBytes); err != nil {
		return err
	}

	if !headBlock.equals(baseHeadBlock) {
		return errors.New("head block doesn't match")
	}

	decoder := NewDecoder(reader)

	return decoder.Decode(v)
}
