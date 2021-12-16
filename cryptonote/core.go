package cryptonote

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Core struct {
	Network *config.Network

	storage Store
	logger  *log.Logger

	// chains contains all the chains belongs to the current node
	chains []*Chain

	// bcLock used for locking blockchain for changes
	bcLock sync.Mutex
}

var (
	ErrAddBlockAlreadyExists              = errors.New("block already exists")
	ErrAddBlockTransactionCountNotMatch   = errors.New("transaction sizes not match")
	ErrAddBlockTransactionSizeMax         = errors.New("transaction size bigger than allowed")
	ErrAddBlockTransactionDeserialization = errors.New("transaction deserialization failed")
	ErrAddBlockRejectedAsOrphaned         = errors.New("rejected as orphaned")
)

func NewCore(network *config.Network, DB Store, logger *log.Logger) (*Core, error) {
	core := &Core{
		Network: network,
		storage: DB,
		logger:  logger,
	}

	if err := core.Init(); err != nil {
		return nil, err
	}

	return core, nil
}

func (c *Core) Init() error {
	if err := c.initDB(); err != nil {
		return errors.New(fmt.Sprintf("failed to initialize storage: %s", err))
	}

	return nil
}

// AddBlock used to add next block to the blockchain
func (c *Core) AddBlock(b *Block, rawTransactions [][]byte) error {
	c.bcLock.Lock()
	defer c.bcLock.Unlock()

	index := b.Index()
	hash, err := b.Hash()
	if err != nil {
		return err
	}

	// c.logger.Tracef("adding block: %d (%s)", index, hash.String())

	hasBlock, err := c.storage.HasBlock(hash)
	if err != nil {
		return err
	}
	if hasBlock {
		c.logger.Errorf("block already exists: %d (%s)", index, hash.String())
		return ErrAddBlockAlreadyExists
	}

	if len(b.TransactionsHashes) != len(rawTransactions) {
		c.logger.Errorf("wrong transaction size: %d (%s)", index, hash.String())
		return ErrAddBlockTransactionCountNotMatch
	}

	// TODO: Check that a block is not orphaned
	if c.findSegmentContainingBlock(hash) == 0 {
		c.logger.Errorf("rejected as orphaned: %d (%s)", index, hash.String())
		return ErrAddBlockRejectedAsOrphaned
	}

	//transactions, transactionsSize, err := c.deserializeTransactions(rawTransactions)
	_, _, err = c.deserializeTransactions(rawTransactions)
	if err != nil {
		return err
	}

	coinbaseTransactionSize, err := b.CoinbaseTransaction.Size()
	if err != nil {
		return err
	}

	if coinbaseTransactionSize > c.Network.MaxTxSize {
		c.logger.Errorf(fmt.Sprintf(
			"coinbase transaction size %d bigger than allowed %d",
			coinbaseTransactionSize,
			c.Network.MaxTxSize,
		))
		return ErrAddBlockTransactionSizeMax
	}

	//blockSize := transactionsSize + coinbaseTransactionSize
	//
	//prevBlockIndex, err := c.BlockHeight(&b.PreviousBlockHash)
	//if err != nil {
	//	c.logger.Errorf("failed to get prev block (%s) index: %s", b.PreviousBlockHash.String(), err)
	//	return err
	//}
	//
	//// TODO: Fetch max height from main blockchain
	//topIndex, err := c.TopIndex()
	//currentBlockchainHeihgt := topIndex + 1
	//if err != nil {
	//	return err
	//}
	// TODO Fetch top index block from current segment
	// bool addOnTop = cache->getTopBlockIndex() == previousBlockIndex;

	if err := c.storage.AppendBlock(b); err != nil {
		return err
	}

	return nil
}

func (c *Core) HasBlock(h *Hash) (bool, error) {
	hasBlock, err := c.storage.HasBlock(h)
	if err != nil {
		return false, nil
	}

	return hasBlock, nil
}

func (c *Core) TopIndex() (uint32, error) {
	height, err := c.storage.TopIndex()
	if err != nil {
		c.logger.Errorf("failed to get height from storage: %s", err)
		return 0, err
	}

	return height, err
}

func (c *Core) BlockHeight(h *Hash) (uint32, error) {
	i, err := c.storage.GetBlockIndexByHash(h)

	if err != nil {
		return 0, err
	}

	return i, nil
}

// BlockByHeight returns block by height
func (c *Core) BlockByHeight(h uint32) (*Block, error) {
	b, err := c.storage.GetBlockByHeight(h)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Core) GenesisBlockHash() (h *Hash, err error) {
	if b, err := c.storage.GetBlockByHeight(0); err == nil {
		if h, err = b.Hash(); err == nil {
			return h, nil
		}
	}

	return nil, fmt.Errorf("failed to get genesis block: %w", err)
}

func (c *Core) TopBlock() (*Block, uint32, error) {
	height, err := c.TopIndex()
	if err != nil {
		return nil, 0, err
	}

	block, err := c.storage.GetBlockByHeight(height)
	if err != nil {
		return nil, 0, err
	}

	return block, height, nil
}

// BuildSparseChain
// IDs pow(2,n) offset, like 2, 4, 8, 16, 32, 64 and so on, and the last one is always genesis block
func (c *Core) BuildSparseChain() ([]Hash, error) {
	var list []Hash

	topBlock, height, err := c.TopBlock()
	if err != nil {
		return nil, err
	}

	topHash, err := topBlock.Hash()
	if err != nil {
		return nil, err
	}

	list = append(list, *topHash)

	for i := uint32(1); i < height; i *= 2 {
		hash, err := c.storage.GetBlockHashByHeight(height - i)
		if err != nil {
			return nil, err
		}

		list = append(list, *hash)
	}

	ghash, err := c.storage.GetBlockHashByHeight(0)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(ghash[:], list[0][:]) && !bytes.Equal(ghash[:], list[len(list)-1][:]) {
		list = append(list, *ghash)
	}

	return list, nil
}

func (c *Core) initDB() error {
	if err := c.storage.Init(c.Network); err != nil {
		return err
	}

	return nil
}

func (c *Core) deserializeTransactions(rawTransactions [][]byte) ([]Transaction, uint64, error) {
	var size uint64
	var transactions []Transaction

	maxTxSize := c.Network.MaxTxSize
	for i, rt := range rawTransactions {
		var t Transaction
		ts := uint64(len(rt))

		if ts > maxTxSize {
			c.logger.Errorf(fmt.Sprintf("transaction size at index %d bigger than allowed %d", i, maxTxSize))
			return nil, 0, ErrAddBlockTransactionSizeMax
		}

		if err := t.Deserialize(bytes.NewReader(rt)); err != nil {
			c.logger.Errorf(fmt.Sprintf("transaction deserialization at index %d failed: %s", i, err))
			return nil, 0, ErrAddBlockTransactionDeserialization
		}

		size += ts
		transactions = append(transactions, t)
	}

	return transactions, size, nil
}

// findChainContainingBlock returns chain containing block by the hash
func (c *Core) findChainContainingBlockHash(*Hash) *Chain {

	return nil
}

// ------------------ Experiments ------------------------------
func (c *Core) findSegmentContainingBlock(h *Hash) int {
	blockSegment := c.findMainChainSegmentContainingBlock(h)

	if blockSegment != 0 {
		return blockSegment
	}

	return c.findMainChainSegmentContainingBlock(h)
}

func (c *Core) findMainChainSegmentContainingBlock(h *Hash) int {
	return 0
}

func (c *Core) findAlternativeSegmentContainingBlock(h *Hash) int {
	return 0
}
