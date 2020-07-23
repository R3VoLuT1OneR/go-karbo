package hash

import (
	"golang.org/x/crypto/sha3"
)

// Keccak hash
func Keccak(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)

	return hash.Sum(nil)
}
