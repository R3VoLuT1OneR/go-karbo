package cryptonote

import (
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
)

// difficultyForNextBlock calculates difficulty for the next block.
func (bc *BlockChain) difficultyForNextBlock(b *Block) (uint64, error) {
	if b.Index() > bc.bestTip.Index() {
		return 0, errors.New(fmt.Sprintf("unknown block index %d, top index is %d", b.Index(), bc.bestTip.Index()))
	}

	nextBlockMajorVersion := bc.Network.GetBlockMajorVersionForHeight(b.Index())
	difficultyBlocksCount := bc.Network.DifficultyBlocksCountByBlockVersion(nextBlockMajorVersion)

	timestamps := bc.lastBlocksTimestamps(difficultyBlocksCount, b)
	cumulativeDifficulties := bc.lastBlocksCumulativeDifficulties(difficultyBlocksCount, b)

	return bc.Network.NextDifficulty(b.Index(), nextBlockMajorVersion, timestamps, cumulativeDifficulties)
}

// TODO: Implement
// blockCumulativeDifficulty returns cumulative difficulty for specific block
func (bc *BlockChain) blockCumulativeDifficulty(b *Block) uint64 {
	return 0
}

// TODO: Refactor make sure it is running fast
func (bc *BlockChain) lastBlocksCumulativeDifficulties(count int, b *Block) []uint64 {
	var difficulties []uint64
	var tempBlock = b

	for count > 0 {
		difficulties = append(difficulties, bc.blockCumulativeDifficulty(tempBlock))

		if tempBlock.PreviousBlockHash == (crypto.Hash{}) {
			break
		}

		tempBlock = bc.getBlockByHash(&tempBlock.PreviousBlockHash)
		count--
	}

	return difficulties
}
