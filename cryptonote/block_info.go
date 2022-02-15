package cryptonote

import "github.com/r3volut1oner/go-karbo/crypto"

// BlockInfo represents additional information about block that is not part of block itself.
//
// Planned use this struct for internal use of chain
type blockInfo struct {
	// Index of the block
	Index uint32

	// Hash of the block
	Hash crypto.Hash

	// CumulativeDifficulty of the POW for the block
	CumulativeDifficulty uint64

	// TotalGeneratedTransactions keeps how many transactions in blockchain including in this block
	TotalGeneratedTransactions uint64

	// TotalGeneratedCoins keeps how many coins generated in blockchain including this block
	TotalGeneratedCoins uint64

	// Timestamp of the block
	Timestamp uint64

	// Size of the block in bytes
	Size uint64
}

// NewBlockInfo return new blockInfo instance
func NewBlockInfo(i uint32, totalGeneratedCoins uint64) *blockInfo {
	return &blockInfo{
		Index:               i,
		TotalGeneratedCoins: totalGeneratedCoins,
	}
}
