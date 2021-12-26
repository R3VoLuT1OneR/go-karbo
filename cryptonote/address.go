package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/crypto/base58"
)

var checksumSize = 4

// Address represents cryptonote address
type Address struct {
	Tag            uint64
	SpendPublicKey crypto.PublicKey
	ViewPublicKey  crypto.PublicKey
	base58         string
}

// NewAddress returns address struct from provided tags
func NewAddress(tag uint64, spendPublicKey, viewPublicKey crypto.PublicKey) (a Address) {
	a.Tag = tag
	a.SpendPublicKey = spendPublicKey
	a.ViewPublicKey = viewPublicKey

	return
}

func (a *Address) Base58() string {
	if a.base58 == "" {
		var b []byte
		b = append(b, a.SpendPublicKey[:]...)
		b = append(b, a.ViewPublicKey[:]...)

		a.base58 = addressEncode(a.Tag, b)
	}

	return a.base58
}

// FromString fill address information from base58 encoded string
func (a *Address) FromString(s string) error {
	tag, data, err := addressDecode(s)

	if err != nil {
		return err
	}

	if len(data) != 64 {
		return errors.New("encoded data has wrong length")
	}

	var spendPublicKeyBytes [32]byte
	var viewPublicKeyBytes [32]byte

	copy(spendPublicKeyBytes[:], data[:32])
	copy(viewPublicKeyBytes[:], data[32:64])

	a.base58 = s
	a.SpendPublicKey = spendPublicKeyBytes
	a.ViewPublicKey = viewPublicKeyBytes
	a.Tag = tag

	return nil
}

// addressDecode decodes base58 encoded address
func addressDecode(addr string) (tag uint64, data []byte, err error) {
	decoded, err := base58.Decode(addr)

	if err != nil {
		return
	}

	if len(decoded) <= checksumSize {
		err = fmt.Errorf("Decoded size is too short %d", len(decoded))
		return
	}

	checksumStart := len(decoded) - checksumSize
	checksum := decoded[checksumStart:]
	ddata := decoded[:checksumStart]

	if !bytes.Equal(checksum, crypto.Keccak(ddata)[:checksumSize]) {
		err = fmt.Errorf("invalid checksum")
		return
	}

	tag, read := binary.Uvarint(decoded[:checksumStart])

	if read <= 0 || read > checksumStart {
		err = fmt.Errorf("Failed read varint")
		return
	}

	data = ddata[read:]
	return
}

// addressEncode encodes base58
func addressEncode(tag uint64, data []byte) string {
	var buf []byte

	// Put first Varin
	vbuf := make([]byte, 16)
	vlen := binary.PutUvarint(vbuf, tag)

	buf = append([]byte{}, vbuf[:vlen]...)
	buf = append(buf, data...)
	buf = append(buf, crypto.Keccak(buf)[:checksumSize]...)

	return base58.Encode(buf)
}
