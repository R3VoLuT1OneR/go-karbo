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
	PushBlock(*Block, *blockInfo, TransactionsDetails) error

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
	getBlockInfoAtIndex(index uint32) *blockInfo
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
