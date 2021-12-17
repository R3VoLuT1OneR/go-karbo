package cryptonote

import (
	"bytes"
	"encoding/hex"
	"github.com/r3volut1oner/go-karbo/config"
)

type BlockChain struct {
	// Network is current network configurations, must stay immutable
	Network *config.Network

	// genesisBlock used for caching genesis block
	genesisBlock *Block
}

// NewBlockChain generates basic blockchain object
func NewBlockChain(network *config.Network) BlockChain {
	bc := BlockChain{
		Network: network,
	}

	return bc
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
