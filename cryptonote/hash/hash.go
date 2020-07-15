package hash

import (
	"golang.org/x/crypto/sha3"
)

// Digest returns digest of provided data
func Digest(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)

	return hash.Sum(nil)
}
