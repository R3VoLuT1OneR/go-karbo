package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	log "github.com/sirupsen/logrus"
	"sort"
)

var (
	ErrCheckpointsAlreadyExists       = errors.New("checkpoints already exists")
	ErrCheckpointsFailed              = errors.New("checkpoints check failed")
	ErrCheckpointsTooDeepReorg        = errors.New("checkpoints too deep reorganisation")
	ErrCheckpointsAltBeforeCheckpoint = errors.New("checkpoints alternative block before checkpoint")
	ErrCheckpointsAltBlockGenesis     = errors.New("checkpoints alternative genesis block not allowed")
)

type Checkpoints interface {
	// AddCheckpoint to the list of all the checkpoints
	AddCheckpoint(height uint32, hash crypto.Hash) error

	// IsInCheckpointZone checks if height in in checkpoints zone
	IsInCheckpointZone(height uint32) bool

	// CheckBlock verifies that hash on provided hash matched the checkpoint
	CheckBlock(height uint32, hash *crypto.Hash) error

	// AlternativeBlockAllowed checks if alternative branch for blockchain allowed on specific height
	// Returns Nil if allowed and an error when is not allowed.
	AlternativeBlockAllowed(bcSize uint32, height uint32) error
}

func NewCheckpoints(logger *log.Logger) Checkpoints {
	return &checkpoints{
		logger: logger,
		points: map[uint32]crypto.Hash{},
	}
}

type checkpoints struct {
	// points represents the checkpoints
	points map[uint32]crypto.Hash

	// logger used for different error messages
	logger *log.Logger

	// pointsSorted keeping indexes sorted, so we can check checkpoint fast
	pointsSorted []uint32
}

func (cp *checkpoints) AddCheckpoint(height uint32, hash crypto.Hash) error {
	if _, ok := cp.points[height]; ok {
		err := ErrCheckpointsAlreadyExists
		cp.prepareLogger(height, &hash).Error(err)
		return err
	}

	cp.points[height] = hash
	cp.pointsSorted = append(cp.pointsSorted, height)

	sort.Slice(cp.pointsSorted, func(i, j int) bool {
		return cp.pointsSorted[i] < cp.pointsSorted[j]
	})

	return nil
}

func (cp *checkpoints) IsInCheckpointZone(height uint32) bool {
	maxHeight := cp.pointsSorted[len(cp.pointsSorted)-1]
	return maxHeight != 0 && height <= maxHeight
}

func (cp *checkpoints) CheckBlock(height uint32, hash *crypto.Hash) error {
	if checkpointHash, ok := cp.points[height]; ok {
		if checkpointHash == *hash {
			return nil
		} else {
			err := ErrCheckpointsFailed
			cp.prepareLogger(height, hash).WithFields(log.Fields{
				"checkpoint_correct_hash": checkpointHash,
			}).Error(err)
			return err
		}
	}

	return nil
}

func (cp *checkpoints) AlternativeBlockAllowed(bcSize uint32, height uint32) error {
	logger := cp.logger.WithFields(log.Fields{
		"checkpoint_blockchain_size": bcSize,
		"checkpoint_height":          height,
	})

	if height <= 1 {
		err := ErrCheckpointsAltBlockGenesis
		logger.Error(err)
		return err
	}

	uw := config.MinedMoneyUnlockWindow
	if height < bcSize-uw && bcSize > uw && !cp.IsInCheckpointZone(height) {
		err := ErrCheckpointsTooDeepReorg
		logger.Error(err)
		return err
	}

	// Is blockchain before first checkpoint?
	if len(cp.pointsSorted) == 0 || bcSize < cp.pointsSorted[0] {
		return nil
	}

	checkpointHeight := cp.pointsSorted[0]
	for _, height := range cp.pointsSorted {
		if height <= bcSize {
			checkpointHeight = height
		}
	}

	if checkpointHeight >= height {
		err := ErrCheckpointsAltBeforeCheckpoint
		logger.Error(err)
		return err
	}

	return nil
}

// prepareLogger adds block details to the logger
func (cp *checkpoints) prepareLogger(height uint32, hash *crypto.Hash) *log.Entry {
	return cp.logger.WithFields(log.Fields{
		"checkpoint_height": height,
		"checkpoint_hash":   hash.String(),
	})
}
