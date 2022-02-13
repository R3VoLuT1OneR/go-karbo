package config

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/crypto"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
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
	AddCheckpoint(index uint32, hash crypto.Hash) error

	// IsInCheckpointZone checks if index in in checkpoints zone
	IsInCheckpointZone(index uint32) bool

	// CheckBlock verifies that hash on provided hash matched the checkpoint
	CheckBlock(index uint32, hash *crypto.Hash) error

	// AlternativeBlockAllowed checks if alternative branch for blockchain allowed on specific index
	// Returns Nil if allowed and an error when is not allowed.
	AlternativeBlockAllowed(bcSize uint32, index uint32) error
}

func NewCheckpoints(logger *log.Logger) Checkpoints {
	return &checkpoints{
		logger:       logger,
		points:       map[uint32]crypto.Hash{},
		pointsSorted: []uint32{},
	}
}

type checkpoints struct {
	// points represents the checkpoints
	points map[uint32]crypto.Hash

	// logger used for different error messages
	logger *log.Logger

	// pointsSorted keeping indexes sorted, so we can check checkpoint fast
	pointsSorted []uint32

	sync.RWMutex
}

func (cp *checkpoints) AddCheckpoint(index uint32, hash crypto.Hash) error {
	cp.Lock()
	defer cp.Unlock()

	if _, ok := cp.points[index]; ok {
		err := ErrCheckpointsAlreadyExists
		cp.prepareLogger(index, &hash).Error(err)
		return err
	}

	cp.points[index] = hash
	cp.pointsSorted = append(cp.pointsSorted, index)

	sort.Slice(cp.pointsSorted, func(i, j int) bool {
		return cp.pointsSorted[i] < cp.pointsSorted[j]
	})

	return nil
}

func (cp *checkpoints) IsInCheckpointZone(index uint32) bool {
	cp.Lock()
	defer cp.Unlock()

	// No checkpoints added
	if len(cp.pointsSorted) == 0 {
		return false
	}

	maxHeight := cp.pointsSorted[len(cp.pointsSorted)-1]
	return maxHeight != 0 && index <= maxHeight
}

func (cp *checkpoints) CheckBlock(index uint32, hash *crypto.Hash) error {
	cp.Lock()
	defer cp.Unlock()

	if checkpointHash, ok := cp.points[index]; ok {
		if checkpointHash == *hash {
			return nil
		} else {
			err := ErrCheckpointsFailed
			cp.prepareLogger(index, hash).WithFields(log.Fields{
				"checkpoint_correct_hash": checkpointHash,
			}).Error(err)
			return err
		}
	}

	return nil
}

func (cp *checkpoints) AlternativeBlockAllowed(bcSize uint32, index uint32) error {
	cp.Lock()
	defer cp.Unlock()

	logger := cp.logger.WithFields(log.Fields{
		"checkpoint_blockchain_size": bcSize,
		"checkpoint_index":           index,
	})

	if index < 1 {
		err := ErrCheckpointsAltBlockGenesis
		logger.Error(err)
		return err
	}

	uw := MinedMoneyUnlockWindow
	if index < bcSize-uw && bcSize > uw && !cp.IsInCheckpointZone(index) {
		err := ErrCheckpointsTooDeepReorg
		logger.Error(err)
		return err
	}

	// Is blockchain before first checkpoint?
	if len(cp.pointsSorted) == 0 || bcSize < cp.pointsSorted[0] {
		return nil
	}

	checkpointHeight := cp.pointsSorted[0]
	for _, index := range cp.pointsSorted {
		if index <= bcSize {
			checkpointHeight = index
		}
	}

	if checkpointHeight >= index {
		err := ErrCheckpointsAltBeforeCheckpoint
		logger.Error(err)
		return err
	}

	return nil
}

// prepareLogger adds block details to the logger
func (cp *checkpoints) prepareLogger(index uint32, hash *crypto.Hash) *log.Entry {
	return cp.logger.WithFields(log.Fields{
		"checkpoint_index": index,
		"checkpoint_hash":  hash.String(),
	})
}
