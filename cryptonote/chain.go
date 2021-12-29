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
	"time"
)

type BlockChain struct {
	// Network is current network configurations, must stay immutable
	Network *config.Network

	// bestTip the higher block in the blockchain
	bestTip *Block

	// tips higher blockchain tips
	tips []*Block

	// genesisBlock used for caching genesis block
	genesisBlock *Block

	blocksIndex map[crypto.Hash]*Block

	// logger for block chain events
	logger *log.Logger

	sync.RWMutex
}

// NewBlockChain generates basic blockchain object
func NewBlockChain(network *config.Network, logger *log.Logger) *BlockChain {
	bc := &BlockChain{
		Network: network,
		logger:  logger,

		tips: []*Block{},
	}

	return bc
}

// AddBlock used for adding new blocks to the blockchain.
//
// It returns nil if block added successfully and ErrAddBlock* in case of error
func (bc *BlockChain) AddBlock(b *Block, rawTransactions [][]byte) error {
	bc.Lock()
	defer bc.Unlock()

	hash := b.Hash()
	height := b.Height()

	blogger := bc.logger.WithFields(log.Fields{
		"block_hash":   hash,
		"block_height": height,
	})

	if bc.haveBlock(hash) {
		err := ErrAddBlockAlreadyExists
		blogger.Error(err)
		return err
	}

	if !bc.haveBlock(&b.PreviousBlockHash) {
		err := ErrAddBlockRejectedAsOrphaned
		blogger.Error(err)
		return err
	}

	coinbaseTransactionSize := b.CoinbaseTransaction.Size()
	if coinbaseTransactionSize > bc.Network.MaxTxSize {
		err := ErrAddBlockTransactionCoinbaseMaxSize
		blogger.Error(err)
		return err
	}

	transactions, transactionsSize, err := bc.deserializeTransactions(blogger, rawTransactions)
	if err != nil {
		return err
	}

	if len(b.TransactionsHashes) != len(transactions) {
		err := ErrAddBlockTransactionCountNotMatch
		blogger.Error(err)
		return err
	}

	prevBlock := bc.getBlockByHash(&b.PreviousBlockHash)
	blockSize := coinbaseTransactionSize + transactionsSize
	if blockSize > bc.Network.MaxBlockSize(uint64(prevBlock.Height())) {
		err := ErrBlockValidationCumulativeSizeTooBig
		blogger.Error(err)
		return err
	}

	if err := bc.validateBlock(blogger, b); err != nil {
		return err
	}

	if b.MajorVersion >= config.BlockMajorVersion5 {
		sigHash := crypto.HashFromBytes(b.HashingBytes())
		outputKey := b.CoinbaseTransaction.Outputs[0].Target.(OutputKey)
		if !b.Signature.Check(&sigHash, &outputKey.PublicKey) {
			err := ErrBlockValidationBlockSignatureMismatch
			blogger.Error(err)
			return err
		}
	}

	currentDifficulty, err := bc.difficultyForNextBlock(prevBlock)

	if err != nil {
		err := ErrAddBlockFailedGetDifficulty
		blogger.Error(err)
		return err
	}

	if currentDifficulty == 0 {
		err := ErrBlockValidationDifficultyOverhead
		blogger.Error(err)
		return err
	}

	// Are we going to add the block to the best blockchain
	addOnTop := bc.bestTip.Height() == prevBlock.Height()

	txAddedHashes := map[crypto.Hash]bool{}
	for i, transaction := range transactions {
		// check if tx hashes in txs blob and header match
		txHash := transaction.Hash()

		btlogger := blogger.WithFields(log.Fields{
			"transaction_index": i,
			"transaction_hash":  txHash.String(),
			"block_hash":        b.TransactionsHashes[i],
		})

		if *txHash != b.TransactionsHashes[i] {
			err := ErrBlockValidationTransactionInconsistency
			btlogger.Error(err)
			return err
		}

		if addOnTop && bc.hasTransaction(txHash) {
			err := ErrBlockValidationDuplicateTransaction
			btlogger.Error(err)
			return err
		}

		// check that there's no duplicate transaction in the block
		if _, ok := txAddedHashes[*txHash]; ok {
			err := ErrBlockValidationDuplicateTransaction
			btlogger.Error(err)
			return err
		}

		txAddedHashes[*txHash] = true

		//fee := uint64(0)
		//minFee := bc.Network.MinimalFee(b.Height())

	}

	return nil
}

// validateBlock validates block.
// Returns and error if block not valid.
func (bc *BlockChain) validateBlock(blogger *log.Entry, b *Block) error {
	if bc.Network.GetBlockMajorVersionForHeight(b.Height()) != b.MajorVersion {
		err := ErrBlockValidationWrongVersion
		blogger.Error(err)
		return err
	}

	if b.MajorVersion == config.BlockMajorVersion2 && b.Parent.MajorVersion > config.BlockMajorVersion1 {
		err := ErrBlockValidationParentBlockWrongVersion
		blogger.WithField("block_parent_major_version", b.Parent.MajorVersion).Error(err)
		return err
	}

	if b.MajorVersion == config.BlockMajorVersion2 || b.MajorVersion == config.BlockMajorVersion3 {
		if len(b.Parent.serialize(false)) > 2048 {
			err := ErrBlockValidationParentBlockSizeTooBig
			blogger.Error(err)
			return err
		}
	}

	if b.Timestamp > uint64(time.Now().Unix())+bc.Network.BlockFutureTimeLimit(b.MajorVersion) {
		err := ErrBlockValidationTimestampTooFarInFuture
		blogger.Error(err)
		return err
	}

	timestampCheckWindow := bc.Network.BlockTimestampCheckWindow(b.MajorVersion)
	lastTimestamps := bc.lastBlocksTimestamps(timestampCheckWindow, b)
	if len(lastTimestamps) >= timestampCheckWindow {
		if b.Timestamp < utils.MedianSlice(lastTimestamps) {
			err := ErrBlockValidationTimestampTooFarInPast
			blogger.Error(err)
			return err
		}
	}

	if len(b.CoinbaseTransaction.Inputs) != 1 {
		err := ErrTransactionInputWrongCount
		blogger.Error(err)
		return err
	}

	if _, ok := b.CoinbaseTransaction.Inputs[0].(InputCoinbase); !ok {
		err := ErrTransactionInputUnexpectedType
		blogger.Error(err)
		return err
	}

	prevBlockHeight := bc.blockHeight(&b.PreviousBlockHash)
	if b.CoinbaseTransaction.Inputs[0].(InputCoinbase).BlockIndex != prevBlockHeight {
		err := ErrTransactionBaseInputWrongBlockIndex
		blogger.Error(err)
		return err
	}

	if uint32(b.CoinbaseTransaction.UnlockHeight) != prevBlockHeight+bc.Network.MinedMoneyUnlockWindow() {
		err := ErrTransactionWrongUnlockTime
		blogger.Error(err)
		return err
	}

	if len(b.CoinbaseTransaction.TransactionSignatures) == 0 {
		err := ErrTransactionBaseInvalidSignaturesCount
		blogger.Error(err)
		return err
	}

	if b.MajorVersion >= config.BlockMajorVersion5 {
		cbTransactionExtraFields, parseError := b.CoinbaseTransaction.ParseExtra()
		if parseError != nil || cbTransactionExtraFields.MiningTag != nil {
			err := ErrBlockValidationBaseTransactionExtraMMTag
			blogger.Error(err)
			return err
		}

		if len(b.CoinbaseTransaction.Outputs) != 1 {
			err := ErrTransactionOutputsInvalidCount
			blogger.Error(err)
			return err
		}

		if _, ok := b.CoinbaseTransaction.Outputs[0].Target.(OutputKey); !ok {
			err := ErrTransactionBaseOutputWrongType
			blogger.Error(err)
			return err
		}
	}

	minerReward := uint64(0)
	for i, output := range b.CoinbaseTransaction.Outputs {
		ologger := blogger.WithField("coinbase_output_index", i)

		if output.Amount == 0 {
			err := ErrTransactionOutputZeroAmount
			ologger.Error(err)
			return err
		}

		switch output.Target.(type) {
		case OutputKey:
			outputKey := output.Target.(OutputKey)
			if !outputKey.Check() {
				err := ErrTransactionOutputInvalidKey
				ologger.Error(err)
				return err
			}
		case OutputMultisignature:
			multisigOutput := output.Target.(OutputMultisignature)
			if int(multisigOutput.RequiredSignaturesCount) > len(multisigOutput.Keys) {
				err := ErrTransactionOutputInvalidRequiredSignaturesCount
				ologger.Error(err)
				return err
			}

			for ki, key := range multisigOutput.Keys {
				if !key.Check() {
					err := ErrTransactionOutputInvalidMultisignatureKey
					ologger.WithField("coinbase_output_key_index", ki).Error(err)
					return err
				}
			}
		default:
			err := ErrTransactionOutputUnknownType
			ologger.Error(err)
			return err
		}

		if minerReward > math.MaxUint64-output.Amount {
			err := ErrTransactionOutputsAmountOverflow
			ologger.Error(err)
			return err
		}

		minerReward += output.Amount
	}

	if b.Height() >= config.UpgradeHeightV4s2 && len(b.CoinbaseTransaction.Extra) > config.MaxExtraSize {
		err := ErrTransactionExtraTooLarge
		blogger.Error(err)
		return err
	}

	return nil
}

// TODO: Implement the method
// blockHeight returns index on the current block
func (bc *BlockChain) blockHeight(h *crypto.Hash) uint32 {
	return 0
}

// TODO: Refactor find better implementation way may use only blockHeight
// topIndex return the index of the top best chain block
func (bc *BlockChain) topIndex() uint32 {
	return bc.bestTip.Height() - 1
}

// TODO: Properly implement this method
// haveBlock return whether the block hash contains in the blockchain
//
// This function is NOT safe for concurrent access
func (bc *BlockChain) haveBlock(h *crypto.Hash) bool {
	_, ok := bc.blocksIndex[*h]
	return ok
}

// HaveBlock return whether the block hash contains in the blockchain
//
// This function is safe for concurrent access.
func (bc *BlockChain) HaveBlock(h *crypto.Hash) bool {
	bc.RLock()
	hasBlock := bc.haveBlock(h)
	bc.RUnlock()
	return hasBlock
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

		if err := bc.genesisBlock.CoinbaseTransaction.Deserialize(reader); err != nil {
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

// TODO: Implement
// getBlockByHash fetch the block from block store
func (bc *BlockChain) getBlockByHash(h *crypto.Hash) *Block {
	return nil
}

// TODO: Implement
// hasTransaction check if transaction is stored in blockchain already
func (bc *BlockChain) hasTransaction(txHash *crypto.Hash) bool {
	return false
}
