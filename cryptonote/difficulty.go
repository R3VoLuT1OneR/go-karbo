package cryptonote

import (
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/utils"
)

// difficultyForNextBlock calculates difficulty for the next block.
func (bc *BlockChain) difficultyForNextBlock(b *Block) (uint64, error) {
	if b.Index() > bc.bestTip.Index() {
		return 0, errors.New(fmt.Sprintf("unknown block hashIndex %d, top hashIndex is %d", b.Index(), bc.bestTip.Index()))
	}

	nextBlockMajorVersion := bc.Network.GetBlockMajorVersion(b.Index())
	difficultyBlocksCount := bc.Network.DifficultyBlocksCountByBlockVersion(nextBlockMajorVersion)

	timestamps := bc.lastBlocksTimestamps(difficultyBlocksCount, b, false)
	cumulativeDifficulties := bc.lastBlocksCumulativeDifficulties(difficultyBlocksCount, b.Index(), false)

	return bc.Network.NextDifficulty(b.Index(), nextBlockMajorVersion, timestamps, cumulativeDifficulties)
}

func (bc *BlockChain) lastBlocksCumulativeDifficulties(count uint32, index uint32, addGenesisBlock bool) []uint64 {
	difficulties := []uint64{}

	tempInfo := bc.storage.getBlockInfoAtIndex(index)
	for i := uint32(1); i <= count; i++ {
		if tempInfo == nil {
			break
		}

		if !addGenesisBlock && tempInfo.Index == 0 {
			break
		}

		difficulties = append(difficulties, tempInfo.CumulativeDifficulty)
		tempInfo = bc.storage.getBlockInfoAtIndex(index - i)
		count--
	}

	return utils.SliceReverse(difficulties)
}
