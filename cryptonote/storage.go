package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto"
)

var (
	ErrStorageNetworkMismatch = errors.New("storage network mismatch")

	ErrStorageBlockExists = errors.New("block exists in storage")
)

// Storage used by blockchain for storing blocks information.
//
// In order to implement new type of storage we need to implement this interface.
type Storage interface {

	// Init method is used for storage initialization for provided network.
	// We may load the chain from persistent storage on this moment.
	//
	// Genesis block must be used for network verification when loading from persistent storage.
	Init(genesisBlock *Block) error

	// TopIndex of current saved blockchain
	TopIndex() (uint32, error)

	// TopBlock returns the best block
	TopBlock() (*Block, error)

	// PushBlock to the blockchain storage.
	PushBlock(*Block, *BlockInfo, TransactionsDetails) error

	// HaveBlock verifies that block is saved in DB
	HaveBlock(*crypto.Hash) bool

	// GetBlock returns block represented by provided hash.
	// Returns nil if block not found.
	GetBlock(*crypto.Hash) *Block

	// HashAtIndex provides block by hash
	// TODO: Review the need for this method. It is used in BuildSparseChain only maybe can be replaced
	//       with some different method there.
	HashAtIndex(uint32) (*crypto.Hash, error)

	// Close database connection
	Close() error

	//// GetBlockIndexByHash returns height for specific block hash
	//GetBlockIndexByHash(*crypto.Hash) (uint32, error)
	//
	//// GetBlockByHeight returns block by height
	//GetBlockByHeight(uint32) (*Block, error)
	//
	//
	//
	//// IsEmpty checks if database is new and empty
	//IsEmpty() (bool, error)
	//

	/**
	 * Methods that must be used only in cryptonote package
	 */

	// getBlockInfoAtIndex return block info at specified index
	getBlockInfoAtIndex(index uint32) *BlockInfo
}

// BlockInfo represents additional information about block that is not part of block itself.
//
// Planned use this struct for internal use of chain
type BlockInfo struct {
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

type TransactionInfo struct {
	// Index of the block
	BlockIndex uint32

	// Index of the transaction in the block
	Index uint32

	// Hash of the transaction
	Hash crypto.Hash

	// UnlockTime of the transaction outputs
	UnlockTime uint64

	// Outputs of the transaction
	Outputs []TransactionOutputTarget

	// GlobalIndexes
	GlobalIndexes []uint32

	// AmountToKeyIndexes global key output indexes spawned in this transaction
	AmountToKeyIndexes map[uint64][]uint32

	// AmountToMultiIndexes global multisignature output indexes spawned in this transaction
	AmountToMultiIndexes map[uint64][]uint32
}

type PackedOutIndex struct {
	// Index of the block
	BlockIndex uint32

	// TransactionIndex index of the transaction in block
	TransactionIndex uint32

	// Index of the output in the transaction
	Index uint32
}

type MultisigAmountGlobalOutputIndexPair struct {
	Amount            uint64
	GlobalOutputIndex uint32
}

// TransactionsDetails used for passing transaction information about block transactions
type TransactionsDetails struct {
	transactions []Transaction

	spentKeyImages []crypto.KeyImage

	spentMultisignatureGlobalIndexes []MultisigAmountGlobalOutputIndexPair
}
