package cryptonote

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/utils"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
)

// BlockChain used for all the blockchain manipulations.
// The storage is responsible for keeping all the blockchain data.
type BlockChain struct {
	// Network is current network configurations, must stay immutable
	Network *config.Network

	// storage used for saving all the blockchain info
	storage Storage

	// logger for block bc events
	logger *log.Logger

	// bestTip the higher block in the blockchain
	bestTip *Block

	// tips higher blockchain tips
	tips []*Block

	blocksIndex map[crypto.Hash]*Block

	// genesisBlock network genesis block.
	genesisBlock *Block

	sync.RWMutex
}

// NewBlockChain generates basic blockchain object
func NewBlockChain(network *config.Network, storage Storage, logger *log.Logger) *BlockChain {
	bc := &BlockChain{
		Network: network,
		logger:  logger,
		storage: storage,

		tips: []*Block{},
	}

	return bc
}

// Init the blockchain including storage
func (bc *BlockChain) Init() error {
	bc.Lock()
	defer bc.Unlock()

	genesisBlock, err := bc.GenesisBlock()
	if err != nil {
		return err
	}

	err = bc.storage.Init(genesisBlock)
	if err != nil {
		return err
	}

	topBlock, err := bc.storage.TopBlock()
	if err != nil {
		return err
	}

	bc.bestTip = topBlock

	return nil
}

// Height returns current blockchain height
func (bc *BlockChain) Height() uint32 {
	// TODO: Implement
	return 0
}

// TopBlock returns current best block
func (bc *BlockChain) TopBlock() *Block {
	bc.RLock()
	topBlock := bc.bestTip
	bc.RUnlock()
	return topBlock
}

// AddBlock used for adding new blocks to the blockchain.
//
// It returns nil if block added successfully and ErrAddBlock* in case of error
func (bc *BlockChain) AddBlock(block *Block, rawTransactions [][]byte) error {
	bc.Lock()
	defer bc.Unlock()

	blockHash := block.Hash()
	blockIndex := block.Index()

	logger := bc.logger.WithFields(log.Fields{
		"block_hash":  blockHash,
		"block_index": blockIndex,
	})

	if bc.haveBlock(blockHash) {
		err := ErrAddBlockAlreadyExists
		logger.Error(err)
		return err
	}

	prevBlock := bc.getBlockByHash(&block.PreviousBlockHash)
	if prevBlock == nil {
		err := ErrAddBlockRejectedAsOrphaned
		logger.Error(err)
		return err
	}

	if blockIndex != prevBlock.Index()+1 {
		logger.WithFields(log.Fields{
			"prev_block_index": prevBlock.Index(),
		})

		err := ErrAddBlockPrevBlockIndexMismatch
		logger.Error(err)
		return err
	}

	coinbaseTransactionSize := block.BaseTransaction.Size()
	if coinbaseTransactionSize > bc.Network.MaxTxSize {
		err := ErrAddBlockTransactionCoinbaseMaxSize
		logger.Error(err)
		return err
	}

	transactions, transactionsSize, err := bc.deserializeTransactions(logger, rawTransactions)
	if err != nil {
		return err
	}

	if len(block.TransactionsHashes) != len(transactions) {
		err := ErrAddBlockTransactionCountNotMatch
		logger.Error(err)
		return err
	}

	blockSize := coinbaseTransactionSize + transactionsSize
	if blockSize > bc.Network.MaxBlockSize(uint64(blockIndex)) {
		err := ErrBlockValidationCumulativeSizeTooBig
		logger.Error(err)
		return err
	}

	minerReward, err := bc.validateBlock(logger, block, prevBlock)
	if err != nil {
		return err
	}

	if block.MajorVersion >= config.BlockMajorVersion5 {
		sigHash := crypto.HashFromBytes(block.HashingBytes())
		outputKey := block.BaseTransaction.Outputs[0].Target.(OutputKey)
		if !block.Signature.Check(&sigHash, &outputKey.PublicKey) {
			err := ErrBlockValidationBlockSignatureMismatch
			logger.Error(err)
			return err
		}
	}

	currentDifficulty, err := bc.difficultyForNextBlock(prevBlock)

	if err != nil {
		err := ErrAddBlockFailedGetDifficulty
		logger.Error(err)
		return err
	}

	if currentDifficulty == 0 {
		err := ErrBlockValidationDifficultyOverhead
		logger.Error(err)
		return err
	}

	// Are we going to add the block to the best blockchain
	addOnTop := bc.bestTip.Index() == prevBlock.Index()

	transactionsValidator := NewBlockTransactionsValidator(bc, blockIndex, logger.Logger)

	txAddedHashes := map[crypto.Hash]bool{}
	for i, transaction := range transactions {
		// check if tx hashes in txs blob and header match
		txHash := transaction.Hash()

		logger := logger.WithFields(log.Fields{
			"transaction_index": i,
			"transaction_hash":  txHash.String(),
			"block_hash":        block.TransactionsHashes[i],
		})

		if *txHash != block.TransactionsHashes[i] {
			err := ErrBlockValidationTransactionInconsistency
			logger.Error(err)
			return err
		}

		if addOnTop && bc.hasTransaction(txHash) {
			err := ErrBlockValidationDuplicateTransaction
			logger.Error(err)
			return err
		}

		// check that there's no duplicate transaction in the block
		if _, ok := txAddedHashes[*txHash]; ok {
			err := ErrBlockValidationDuplicateTransaction
			logger.Error(err)
			return err
		}

		txAddedHashes[*txHash] = true

		if err := transactionsValidator.validate(&transaction); err != nil {
			// TODO: Remove transaction from memory pool
			return err
		}
	}

	alreadyGeneratedCoins := bc.getAlreadyGeneratedCoins(prevBlock.Index())
	lastBlockSizes := bc.GetLastBlocksSizes(prevBlock.Index(), true)
	blockSizeMedian := utils.MedianSlice(lastBlockSizes)

	expectedReward, emissionChange, err := bc.Network.GetBlockReward(
		block.MajorVersion, blockSizeMedian, blockSize, alreadyGeneratedCoins, transactionsValidator.cumulativeFee,
	)

	if err != nil {
		err := ErrBlockValidationCumulativeSizeTooBig
		logger.Error(err)
		return err
	}

	if expectedReward != minerReward {
		logger := logger.WithFields(log.Fields{
			"block_expected_reward": expectedReward,
			"block_miner_reward":    minerReward,
		})

		err := ErrBlockValidationBlockRewardMismatch
		logger.Error(err)
		return err
	}

	if bc.Network.Checkpoints.IsInCheckpointZone(blockIndex) {
		if err := bc.Network.Checkpoints.CheckBlock(blockIndex, block.Hash()); err != nil {
			err := ErrBlockValidationCheckpointBlockHashMismatch
			logger.Error(err)
			return err
		}
	} else {
		if err := bc.checkProofOfWork(block, currentDifficulty); err != nil {
			err := ErrBlockValidationProofOfWorkTooWeak
			logger.Error(err)
			return err
		}
	}

	// TODO: Remove
	logger.Info(emissionChange)

	return nil
}

// BuildSparseChain
// IDs pow(2,n) offset, like 2, 4, 8, 16, 32, 64 and so on, and the last one is always genesis block
func (bc *BlockChain) BuildSparseChain() ([]crypto.Hash, error) {
	var list []crypto.Hash

	topBlock := bc.TopBlock()
	topHash := topBlock.Hash()

	list = append(list, *topHash)

	height := bc.Height()
	for i := uint32(1); i < height; i *= 2 {
		hash, err := bc.storage.HashAtIndex(height - i)
		if err != nil {
			return nil, err
		}

		list = append(list, *hash)
	}

	genesisHash, err := bc.storage.HashAtIndex(0)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(genesisHash[:], list[0][:]) && !bytes.Equal(genesisHash[:], list[len(list)-1][:]) {
		list = append(list, *genesisHash)
	}

	return list, nil
}

// HaveBlock returns whether the chain have block represented by provided hash.
//
// This function is safe for concurrent access.
func (bc *BlockChain) HaveBlock(hash *crypto.Hash) bool {
	bc.RLock()
	have := bc.haveBlock(hash)
	bc.RUnlock()
	return have
}

// validateBlock validates block.
// Returns and error if block not valid.
func (bc *BlockChain) validateBlock(blogger *log.Entry, block *Block, prevBlock *Block) (uint64, error) {
	if bc.Network.GetBlockMajorVersionForHeight(block.Index()) != block.MajorVersion {
		err := ErrBlockValidationWrongVersion
		blogger.Error(err)
		return 0, err
	}

	if block.MajorVersion == config.BlockMajorVersion2 && block.Parent.MajorVersion > config.BlockMajorVersion1 {
		err := ErrBlockValidationParentBlockWrongVersion
		blogger.WithField("block_parent_major_version", block.Parent.MajorVersion).Error(err)
		return 0, err
	}

	if block.MajorVersion == config.BlockMajorVersion2 || block.MajorVersion == config.BlockMajorVersion3 {
		if len(block.Parent.serialize(false)) > 2048 {
			err := ErrBlockValidationParentBlockSizeTooBig
			blogger.Error(err)
			return 0, err
		}
	}

	if block.Timestamp > bc.Network.Timestamp()+bc.Network.BlockFutureTimeLimit(block.MajorVersion) {
		err := ErrBlockValidationTimestampTooFarInFuture
		blogger.Error(err)
		return 0, err
	}

	timestampCheckWindow := bc.Network.BlockTimestampCheckWindow(block.MajorVersion)
	lastTimestamps := bc.lastBlocksTimestamps(timestampCheckWindow, block)
	if len(lastTimestamps) >= timestampCheckWindow {
		if block.Timestamp < utils.MedianSlice(lastTimestamps) {
			err := ErrBlockValidationTimestampTooFarInPast
			blogger.Error(err)
			return 0, err
		}
	}

	if len(block.BaseTransaction.Inputs) != 1 {
		err := ErrTransactionInputWrongCount
		blogger.Error(err)
		return 0, err
	}

	if _, ok := block.BaseTransaction.Inputs[0].(InputCoinbase); !ok {
		err := ErrTransactionInputUnexpectedType
		blogger.Error(err)
		return 0, err
	}

	prevBlockHeight := prevBlock.Index()
	if block.BaseTransaction.Inputs[0].(InputCoinbase).BlockIndex != prevBlockHeight {
		err := ErrTransactionBaseInputWrongBlockIndex
		blogger.Error(err)
		return 0, err
	}

	if uint32(block.BaseTransaction.UnlockHeight) != prevBlockHeight+bc.Network.MinedMoneyUnlockWindow() {
		err := ErrTransactionWrongUnlockTime
		blogger.Error(err)
		return 0, err
	}

	if len(block.BaseTransaction.TransactionSignatures) == 0 {
		err := ErrTransactionBaseInvalidSignaturesCount
		blogger.Error(err)
		return 0, err
	}

	if block.MajorVersion >= config.BlockMajorVersion5 {
		cbTransactionExtraFields, parseError := block.BaseTransaction.ParseExtra()
		if parseError != nil || cbTransactionExtraFields.MiningTag != nil {
			err := ErrBlockValidationBaseTransactionExtraMMTag
			blogger.Error(err)
			return 0, err
		}

		if len(block.BaseTransaction.Outputs) != 1 {
			err := ErrTransactionOutputsInvalidCount
			blogger.Error(err)
			return 0, err
		}

		if _, ok := block.BaseTransaction.Outputs[0].Target.(OutputKey); !ok {
			err := ErrTransactionBaseOutputWrongType
			blogger.Error(err)
			return 0, err
		}
	}

	minerReward := uint64(0)
	for i, output := range block.BaseTransaction.Outputs {
		ologger := blogger.WithField("coinbase_output_index", i)

		if output.Amount == 0 {
			err := ErrTransactionOutputZeroAmount
			ologger.Error(err)
			return 0, err
		}

		switch output.Target.(type) {
		case OutputKey:
			outputKey := output.Target.(OutputKey)
			if !outputKey.Check() {
				err := ErrTransactionOutputInvalidKey
				ologger.Error(err)
				return 0, err
			}
		case OutputMultisignature:
			multisigOutput := output.Target.(OutputMultisignature)
			if int(multisigOutput.RequiredSignaturesCount) > len(multisigOutput.Keys) {
				err := ErrTransactionOutputInvalidRequiredSignaturesCount
				ologger.Error(err)
				return 0, err
			}

			for ki, key := range multisigOutput.Keys {
				if !key.Check() {
					err := ErrTransactionOutputInvalidMultisignatureKey
					ologger.WithField("coinbase_output_key_index", ki).Error(err)
					return 0, err
				}
			}
		default:
			err := ErrTransactionOutputUnknownType
			ologger.Error(err)
			return 0, err
		}

		if minerReward > math.MaxUint64-output.Amount {
			err := ErrTransactionOutputsAmountOverflow
			ologger.Error(err)
			return 0, err
		}

		minerReward += output.Amount
	}

	if block.Index() >= config.UpgradeHeightV4s2 && len(block.BaseTransaction.Extra) > config.MaxExtraSize {
		err := ErrTransactionExtraTooLarge
		blogger.Error(err)
		return 0, err
	}

	return minerReward, nil
}

// TODO: Properly implement this method
// haveBlock return whether the block hash contains in the blockchain
//
// This function is NOT safe for concurrent access
func (bc *BlockChain) haveBlock(h *crypto.Hash) bool {
	if bc.storage.HaveBlock(h) {
		return true
	}

	// TODO: Implement search for a block in side chains
	return false
}

// GenesisBlock returns first basic block of the blockchain
func (bc *BlockChain) GenesisBlock() (*Block, error) {
	if bc.genesisBlock == nil {
		bc.genesisBlock = &Block{}
		genesisTransactionBytes, err := hex.DecodeString(bc.Network.GenesisCoinbaseTxHex)
		reader := bytes.NewReader(genesisTransactionBytes)

		if err != nil {
			return nil, err
		}

		if err := bc.genesisBlock.BaseTransaction.Deserialize(reader); err != nil {
			return nil, err
		}

		bc.genesisBlock.MajorVersion = config.BlockMajorVersion1
		bc.genesisBlock.MinorVersion = config.BlockMinorVersion0
		bc.genesisBlock.Timestamp = bc.Network.GenesisTimestamp
		bc.genesisBlock.Nonce = bc.Network.GenesisNonce
	}

	return bc.genesisBlock, nil
}

// deserializeTransactions deserializes transactions to object, transactions are passing basic data validation.
// Function write to the log on error, so no need to log the error on the caller side.
func (bc *BlockChain) deserializeTransactions(blogger *log.Entry, rt [][]byte) ([]Transaction, uint64, error) {
	transactions := make([]Transaction, len(rt))
	transactionsSize := uint64(0)

	for i, t := range transactions {
		tsSize := uint64(len(rt[i]))
		tsLogger := blogger.WithFields(log.Fields{
			"transaction_size":  tsSize,
			"transaction_index": i,
		})

		if tsSize > bc.Network.MaxTxSize {
			err := ErrAddBlockTransactionSizeMax
			tsLogger.Error(err)
			return nil, 0, err
		}

		r := bytes.NewReader(rt[i])
		if err := t.Deserialize(r); err != nil {
			tsLogger.Error(fmt.Errorf("%s: %w", ErrAddBlockTransactionDeserialization.Error(), err))
			return nil, 0, ErrAddBlockTransactionDeserialization
		}

		transactionsSize += tsSize
	}

	return transactions, transactionsSize, nil
}

// lastBlocksTimestamps fetches the timestamps of the
func (bc *BlockChain) lastBlocksTimestamps(count int, b *Block) []uint64 {
	var timestamps []uint64
	var tempBlock = b

	for count > 0 {
		timestamps = append(timestamps, tempBlock.Timestamp)

		if tempBlock.PreviousBlockHash == (crypto.Hash{}) {
			break
		}

		tempBlock = bc.getBlockByHash(&tempBlock.PreviousBlockHash)
		count--
	}

	return timestamps
}

// IsTransactionSpendTimeUnlocked check
func (bc *BlockChain) IsTransactionSpendTimeUnlocked(unlockTime uint64, blockHeight uint32) bool {
	// Interpret as block height
	if unlockTime < bc.Network.MaxBlockNumber {
		return uint64(blockHeight)+config.LockedTxAllowedDeltaBlocks >= unlockTime
	}

	// Interpret as timestamp
	return bc.Network.Timestamp()+config.LockedTxAllowedDeltaSecond >= unlockTime
}

func (bc *BlockChain) DecomposeAmount(amount uint64, dustThreshold uint64) []uint64 {
	chunks, dusts := DecomposeAmountIntoDigits(amount, dustThreshold)

	return append(chunks, dusts...)
}

// checkProofOfWork verify block proof of work
// TODO: Implement
func (bc *BlockChain) checkProofOfWork(block *Block, difficulty uint64) error {
	return nil
}

// getAlreadyGeneratedCoins returns generated coins on specified height
// TODO: Implement
func (bc *BlockChain) getAlreadyGeneratedCoins(height uint32) uint64 {
	return uint64(0)
}

// GetLastBlocksSizes returns last block sizes
// TODO: Implement
func (bc *BlockChain) GetLastBlocksSizes(height uint32, useGenesis bool) []uint64 {
	return nil
}

// IsSpent
// TODO: Implement
func (bc *BlockChain) IsSpent(image crypto.KeyImage, height uint32) bool {
	return false
}

// ExtractKeyOutputKeys
// TODO: Implement
func (bc *BlockChain) ExtractKeyOutputKeys(amount uint64, height uint32, globalIndexes []uint32) ([]crypto.PublicKey, error) {
	return nil, nil
}

// getBlockByHash fetch the block from block store
// TODO: Implement
func (bc *BlockChain) getBlockByHash(h *crypto.Hash) *Block {
	return nil
}

// hasTransaction check if transaction is stored in blockchain already
// TODO: Implement
func (bc *BlockChain) hasTransaction(txHash *crypto.Hash) bool {
	return false
}

// IsMultiSignatureOutputExists check if multisig output exists
// TODO: Implement
func (bc *BlockChain) IsMultiSignatureOutputExists(amount uint64, globalIndex uint32, blockHeight uint32) (*OutputMultisignature, uint64, bool) {
	return nil, 0, false
}

// IsMultiSignatureSpent check
// TODO: Implement
func (bc *BlockChain) IsMultiSignatureSpent(amount uint64, globalIndex uint32, blockHeight uint32) bool {
	return false
}
