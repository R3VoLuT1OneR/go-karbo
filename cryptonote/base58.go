// Package base58 is implementation of CryptoNote Base58.
//
// CryptoNote Base58 is a binary-to-text encoding scheme used to
// represent arbitrary binary data as a sequence of alphanumeric
// characters.
//
// CryptoNote Base58 uses the following alphabet:
//
// 	 123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz
//
// The input is split into 8-byte blocks. The last block may be smaller
// than 8 bytes. Each block is interpreted as a big-endian integer,
// converted into base 58 (again, big-endian), and encoded using the
// alphabet shown above. The number of base-58 digits used to encode a
// block is the smallest number of digits sufficient to encode every
// block of the same size. For example, 8-byte blocks are encoded using
// 11 characters.
package cryptonote

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto/hash"
)

const (
	// alphabet is the modified base58 alphabet used by Bitcoin.
	alphabetSize = 58

	fullBlockSize        = 8
	fullEncodedBlockSize = 11

	checksumSize = 4
)

var (
	blockSizes        = []int{0, -1, 1, 2, -1, 3, 4, 5, -1, 6, 7, 8}
	encodedBlockSizes = []int{0, 2, 3, 5, 6, 7, 9, 10, 11}

	alphabet       = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	alphabetRevert = map[rune]uint64{
		'1': 0, '2': 1, '3': 2, '4': 3, '5': 4, '6': 5, '7': 6, '8': 7,
		'9': 8, 'A': 9, 'B': 10, 'C': 11, 'D': 12, 'E': 13, 'F': 14, 'G': 15,
		'H': 16, 'J': 17, 'K': 18, 'L': 19, 'M': 20, 'N': 21, 'P': 22, 'Q': 23,
		'R': 24, 'S': 25, 'T': 26, 'U': 27, 'V': 28, 'W': 29, 'X': 30, 'Y': 31,
		'Z': 32, 'a': 33, 'b': 34, 'c': 35, 'd': 36, 'e': 37, 'f': 38, 'g': 39,
		'h': 40, 'i': 41, 'j': 42, 'k': 43, 'm': 44, 'n': 45, 'o': 46, 'p': 47,
		'q': 48, 'r': 49, 's': 50, 't': 51, 'u': 52, 'v': 53, 'w': 54, 'x': 55,
		'y': 56, 'z': 57,
	}
)

func uint64touint8be(u uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, u)

	var zeroSize int
	for i, b := range bytes {
		if b != byte(0x00) || i == 7 {
			break
		}

		zeroSize++
	}

	return bytes[zeroSize:]
}

func decodeBlock(b string) ([]byte, error) {
	answer := uint64(0)
	order := uint64(1)
	resSize := blockSizes[len(b)]

	if len(b) >= len(blockSizes) || resSize < 0 {
		return nil, fmt.Errorf("invalid block size")
	}

	for i := len(b) - 1; i >= 0; i-- {
		num, found := alphabetRevert[rune(b[i])]

		if !found {
			return nil, fmt.Errorf("invalid character '%c' found", b[i])
		}

		mul := num * order
		if mul/order != num {
			return nil, fmt.Errorf("Block overflow")
		}

		tmp := answer + mul
		if tmp < answer {
			return nil, fmt.Errorf("Block overflow")
		}

		answer = tmp
		order *= alphabetSize
	}

	if resSize < fullBlockSize && uint64(1)<<(8*resSize) <= answer {
		return nil, fmt.Errorf("Block overflow")
	}

	be := uint64touint8be(answer)
	missingBytes := resSize - len(be)

	if missingBytes < 0 {
		return nil, fmt.Errorf("Block '%s' size overflow", b)
	}

	return append(make([]byte, missingBytes), be...), nil
}

func decode(s string) (result []byte, err error) {
	fullBlockCount := len(s) / fullEncodedBlockSize
	lastBlockSize := len(s) % fullEncodedBlockSize

	for i := 0; i < fullBlockCount; i++ {
		res, err := decodeBlock(s[i*11 : (i+1)*11])
		if err != nil {
			return nil, err
		}

		result = append(result, res...)
	}

	if lastBlockSize > 0 {
		res, err := decodeBlock(s[fullBlockCount*11:])
		if err != nil {
			return nil, err
		}

		result = append(result, res...)
	}

	return
}

// DecodeAddr decodes base58 encoded address
func DecodeAddr(addr string) (tag uint64, data []byte, err error) {
	decoded, err := decode(addr)

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

	if !bytes.Equal(checksum, hash.Keccak(&ddata)[:checksumSize]) {
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

func uint8beTo64(b []byte) uint64 {
	emptyBytes := make([]byte, 8-len(b))
	uint64bytes := append(emptyBytes, b...)

	return binary.BigEndian.Uint64(uint64bytes)
}

func encodeBlock(src []byte) string {
	var num uint64
	var i int

	if len(src) > fullBlockSize {
		num = uint8beTo64(src[:fullBlockSize])
		i = fullEncodedBlockSize
	} else {
		num = uint8beTo64(src)
		i = encodedBlockSizes[len(src)]
	}

	output := make([]byte, i)
	for i--; i >= 0; i-- {
		output[i] = alphabet[num%alphabetSize]
		num /= alphabetSize
	}

	return string(output)
}

func encode(src []byte) (result string) {
	if len(src) == 0 {
		return
	}

	fullBlockCount := len(src) / fullBlockSize
	lastBlockSize := len(src) % fullBlockSize

	for i := 0; i < fullBlockCount; i++ {
		result += encodeBlock(src[i*fullBlockSize:])
	}

	if lastBlockSize > 0 {
		result += encodeBlock(src[fullBlockCount*fullBlockSize:])
	}

	return
}

// EncodeAddr encodes base58
func EncodeAddr(tag uint64, data []byte) string {
	var buf []byte

	// Put first Varin
	vbuf := make([]byte, 16)
	vlen := binary.PutUvarint(vbuf, tag)

	buf = append([]byte{}, vbuf[:vlen]...)
	buf = append(buf, data...)
	buf = append(buf, hash.Keccak(&buf)[:checksumSize]...)

	return encode(buf)
}
