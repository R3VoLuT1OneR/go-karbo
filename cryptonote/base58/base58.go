package base58

import (
	"math/big"
)

const (
	// alphabet is the modified base58 alphabet used by Bitcoin.
	alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	alphabetIdx0 = '1'
)

var b58 = map[string]int64{
	"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7,
	"9": 8, "A": 9, "B": 10, "C": 11, "D": 12, "E": 13, "F": 14, "G": 15,
	"H": 16, "J": 17, "K": 18, "L": 19, "M": 20, "N": 21, "P": 22, "Q": 23,
	"R": 24, "S": 25, "T": 26, "U": 27, "V": 28, "W": 29, "X": 30, "Y": 31,
	"Z": 32, "a": 33, "b": 34, "c": 35, "d": 36, "e": 37, "f": 38, "g": 39,
	"h": 40, "i": 41, "j": 42, "k": 43, "m": 44, "n": 45, "o": 46, "p": 47,
	"q": 48, "r": 49, "s": 50, "t": 51, "u": 52, "v": 53, "w": 54, "x": 55,
	"y": 56, "z": 57,
}

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

func decodeBlock(b string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	scratch := new(big.Int)
	for i := len(b) - 1; i >= 0; i-- {
		tmp, exists := b58[string(b[i])]
		if !exists {
			return []byte("")
		}

		scratch.SetInt64(int64(tmp))
		scratch.Mul(j, scratch)
		answer.Add(answer, scratch)
		j.Mul(j, bigRadix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabetIdx0 {
			break
		}
	}

	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}

func encodeBlock(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}

// Decode cryptonote 58 to raw
func Decode(s string) (result []byte) {
	fullblocks := len(s) / 11
	for i := 0; i < fullblocks; i++ {
		result = append(result, decodeBlock(s[i*11:(i+1)*11])...)
	}
	if len(s)%11 > 0 {
		result = append(result, decodeBlock(s[fullblocks*11:])...)
	}
	return
}

// Encode raw to cryptonote base58 string
func Encode(b []byte) (result string) {
	return
}
