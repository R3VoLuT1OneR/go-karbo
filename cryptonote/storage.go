package cryptonote

// Storage used by blockchain for storing blocks information.
//
// In order to implement new type of storage we need to implement this interface.
type Storage interface {
	//// GetBlockHashByHeight provides block by hash
	//GetBlockHashByHeight(uint32) (*crypto.Hash, error)
	//
	//// GetBlockIndexByHash returns height for specific block hash
	//GetBlockIndexByHash(*crypto.Hash) (uint32, error)
	//
	//// GetBlockByHeight returns block by height
	//GetBlockByHeight(uint32) (*Block, error)
	//
	//// AppendBlock to database persistence layer.
	//AppendBlock(*Block) error
	//
	//// HasBlock verifies that block is saved in DB
	//HasBlock(*crypto.Hash) (bool, error)
	//
	//// TopIndex of current saved blockchain
	//TopIndex() (uint32, error)
	//
	//// IsEmpty checks if database is new and empty
	//IsEmpty() (bool, error)
	//
	//// Close database connection
	//Close() error
}
