package cryptonote

//type Blockchain struct {
//	BlockChain *BlockChain
//
//	storage Storage
//	logger  *log.Logger
//
//	// bcLock used for locking blockchain for read and write
//	bcLock sync.RWMutex
//}
//
//func NewCore(bc *BlockChain, DB Storage, logger *log.Logger) (*Blockchain, error) {
//	core := &Blockchain{
//		BlockChain: bc,
//		storage:    DB,
//		logger:     logger,
//	}
//
//	if err := core.Init(); err != nil {
//		return nil, err
//	}
//
//	return core, nil
//}
//
//func (c *Blockchain) Init() error {
//	if err := c.initDB(); err != nil {
//		return errors.New(fmt.Sprintf("failed to initialize storage: %s", err))
//	}
//
//	return nil
//}
//
//// AddBlock used to add next block to the blockchain
//func (c *Blockchain) AddBlock(b *Block, rawTransactions [][]byte) error {
//	c.bcLock.Lock()
//	defer c.bcLock.Unlock()
//
//	hashIndex := b.Index()
//	hash := b.Hash()
//	blockStr := fmt.Sprintf("%s (%d)", hash, hashIndex)
//
//	c.logger.Tracef("adding block: %s", blockStr)
//
//	// Verify that block is not added to blockchain
//	hasBlock, err := c.storage.HaveBlock(hash)
//	if err != nil {
//		return fmt.Errorf("unexpected error: %w", err)
//	}
//	if hasBlock {
//		c.logger.Errorf("block already exists: %s", blockStr)
//		return ErrAddBlockAlreadyExists
//	}
//
//	// Number of block transaction hashes must be sane as raw transactions that will be saved
//	if len(b.TransactionsHashes) != len(rawTransactions) {
//		c.logger.Errorf("wrong transaction size: %s", blockStr)
//		return ErrAddBlockTransactionCountNotMatch
//	}
//
//	// TODO: Check that a block is not orphaned
//	if c.findSegmentContainingBlock(hash) == 0 {
//		c.logger.Errorf("rejected as orphaned: %d (%s)", hashIndex, hash.String())
//		return ErrAddBlockRejectedAsOrphaned
//	}
//
//	//transactions, transactionsSize, err := c.deserializeTransactions(rawTransactions)
//	_, _, err = c.deserializeTransactions(rawTransactions)
//	if err != nil {
//		return err
//	}
//
//	coinbaseTransactionSize := b.CoinbaseTransaction.Size()
//	if coinbaseTransactionSize > c.BlockChain.Network.MaxTxSize {
//		c.logger.Errorf(fmt.Sprintf(
//			"coinbase transaction size %d bigger than allowed %d",
//			coinbaseTransactionSize,
//			c.BlockChain.Network.MaxTxSize,
//		))
//		return ErrAddBlockTransactionSizeMax
//	}
//
//	//blockSize := transactionsSize + coinbaseTransactionSize
//	//
//	//prevBlockIndex, err := c.BlockHeight(&b.PreviousBlockHash)
//	//if err != nil {
//	//	c.logger.Errorf("failed to get prev block (%s) hashIndex: %s", b.PreviousBlockHash.String(), err)
//	//	return err
//	//}
//	//
//	//// TODO: Fetch max height from main blockchain
//	//topIndex, err := c.TopIndex()
//	//currentBlockchainHeihgt := topIndex + 1
//	//if err != nil {
//	//	return err
//	//}
//	// TODO Fetch top hashIndex block from current segment
//	// bool addOnTop = cache->getTopBlockIndex() == previousBlockIndex;
//
//	if err := c.storage.AppendBlock(b); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (c *Blockchain) HaveBlock(h *crypto.Hash) (bool, error) {
//	hasBlock, err := c.storage.HaveBlock(h)
//	if err != nil {
//		return false, nil
//	}
//
//	return hasBlock, nil
//}
//
//func (c *Blockchain) TopIndex() (uint32, error) {
//	height, err := c.storage.TopIndex()
//	if err != nil {
//		c.logger.Errorf("failed to get height from storage: %s", err)
//		return 0, err
//	}
//
//	return height, err
//}
//
//func (c *Blockchain) BlockHeight(h *crypto.Hash) (uint32, error) {
//	i, err := c.storage.GetBlockIndexByHash(h)
//
//	if err != nil {
//		return 0, err
//	}
//
//	return i, nil
//}
//
//// BlockByHeight returns block by height
//func (c *Blockchain) BlockByHeight(h uint32) (*Block, error) {
//	b, err := c.storage.GetBlockByHeight(h)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return b, nil
//}
//
//func (c *Blockchain) GenesisBlockHash() (h *crypto.Hash, err error) {
//	if b, err := c.storage.GetBlockByHeight(0); err == nil {
//		return b.Hash(), nil
//	}
//
//	return nil, fmt.Errorf("failed to get genesis block: %w", err)
//}
//
//func (c *Blockchain) TopBlock() (*Block, uint32, error) {
//	height, err := c.TopIndex()
//	if err != nil {
//		return nil, 0, err
//	}
//
//	block, err := c.storage.GetBlockByHeight(height)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return block, height, nil
//}

//// BuildSparseChain
//// IDs pow(2,n) offset, like 2, 4, 8, 16, 32, 64 and so on, and the last one is always genesis block
//func (c *Blockchain) BuildSparseChain() ([]crypto.Hash, error) {
//	var list []crypto.Hash
//
//	bestTip, height, err := c.TopBlock()
//	if err != nil {
//		return nil, err
//	}
//
//	topHash := bestTip.Hash()
//
//	list = append(list, *topHash)
//
//	for i := uint32(1); i < height; i *= 2 {
//		hash, err := c.storage.BlockAtIndex(height - i)
//		if err != nil {
//			return nil, err
//		}
//
//		list = append(list, *hash)
//	}
//
//	ghash, err := c.storage.BlockAtIndex(0)
//	if err != nil {
//		return nil, err
//	}
//
//	if !bytes.Equal(ghash[:], list[0][:]) && !bytes.Equal(ghash[:], list[len(list)-1][:]) {
//		list = append(list, *ghash)
//	}
//
//	return list, nil
//}
//
//func (c *Blockchain) initDB() error {
//	if err := c.storage.Init(c.BlockChain); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (c *Blockchain) deserializeTransactions(rawTransactions [][]byte) ([]Transaction, uint64, error) {
//	var size uint64
//	var transactions []Transaction
//
//	maxTxSize := c.BlockChain.Network.MaxTxSize
//	for i, rt := range rawTransactions {
//		var t Transaction
//		ts := uint64(len(rt))
//
//		if ts > maxTxSize {
//			c.logger.Errorf(fmt.Sprintf("transaction size at hashIndex %d bigger than allowed %d", i, maxTxSize))
//			return nil, 0, ErrAddBlockTransactionSizeMax
//		}
//
//		if err := t.Deserialize(bytes.NewReader(rt)); err != nil {
//			c.logger.Errorf(fmt.Sprintf("transaction deserialization at hashIndex %d failed: %s", i, err))
//			return nil, 0, ErrAddBlockTransactionDeserialization
//		}
//
//		size += ts
//		transactions = append(transactions, t)
//	}
//
//	return transactions, size, nil
//}
//
//// ------------------ Experiments ------------------------------
//func (c *Blockchain) findSegmentContainingBlock(h *crypto.Hash) int {
//	blockSegment := c.findMainChainSegmentContainingBlock(h)
//
//	if blockSegment != 0 {
//		return blockSegment
//	}
//
//	return c.findMainChainSegmentContainingBlock(h)
//}
//
//func (c *Blockchain) findMainChainSegmentContainingBlock(h *crypto.Hash) int {
//	return 0
//}
//
//func (c *Blockchain) findAlternativeSegmentContainingBlock(h *crypto.Hash) int {
//	return 0
//}
