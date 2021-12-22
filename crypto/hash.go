package crypto

import "golang.org/x/crypto/sha3"

// Keccak hash function for fast hashing
// it is the "Crypto::cn_fast_hash" method in C++ implementation
func Keccak(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()

	if _, err := hash.Write(data); err != nil {
		panic(err)
	}

	return hash.Sum(nil)
}
