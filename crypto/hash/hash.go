package hash

import (
	"golang.org/x/crypto/sha3"
)

// Keccak hash
func Keccak(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()

	if _, err := hash.Write(data); err != nil {
		panic(err)
	}

	return hash.Sum(nil)
}
