package config

import (
	"encoding/hex"
	"github.com/r3volut1oner/go-karbo/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func loggerDevNull() *log.Logger {
	logger := log.New()
	logger.Out = ioutil.Discard

	return logger
}

func TestCheckpoints_HandlesEmptyCheckpoints(t *testing.T) {
	cp := NewCheckpoints(loggerDevNull())

	assert.Equal(t, ErrCheckpointsAltBlockGenesis, cp.AlternativeBlockAllowed(0, 0))

	assert.Nil(t, cp.AlternativeBlockAllowed(2, 2))
	assert.Nil(t, cp.AlternativeBlockAllowed(2, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(9, 2))
}

func TestCheckpoints_HandlesOneCheckpoint(t *testing.T) {
	hashBytes, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	cp := NewCheckpoints(loggerDevNull())
	_ = cp.AddCheckpoint(5, crypto.HashFromBytes(hashBytes))

	assert.Equal(t, ErrCheckpointsAltBlockGenesis, cp.AlternativeBlockAllowed(0, 0))

	assert.Nil(t, cp.AlternativeBlockAllowed(1, 2))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 4))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 9))

	assert.Nil(t, cp.AlternativeBlockAllowed(4, 2))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 4))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 9))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 9))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 9))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(9, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(9, 9))
}

func TestCheckpoints_HandlesTwoCheckpoints(t *testing.T) {

	hashBytes, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

	cp := NewCheckpoints(loggerDevNull())
	_ = cp.AddCheckpoint(5, crypto.HashFromBytes(hashBytes))
	_ = cp.AddCheckpoint(9, crypto.HashFromBytes(hashBytes))

	assert.Nil(t, cp.AlternativeBlockAllowed(1, 2))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 4))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 8))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(1, 11))

	assert.Nil(t, cp.AlternativeBlockAllowed(4, 2))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 4))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 8))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(4, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(5, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 8))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(5, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(6, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 8))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(6, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(8, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(8, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(8, 5))
	assert.Nil(t, cp.AlternativeBlockAllowed(8, 6))
	assert.Nil(t, cp.AlternativeBlockAllowed(8, 8))
	assert.Nil(t, cp.AlternativeBlockAllowed(8, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(8, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(8, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 5))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 6))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 8))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(9, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(9, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(9, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 5))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 6))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 8))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(10, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(10, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(10, 11))

	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 2))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 4))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 5))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 6))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 8))
	assert.Equal(t, ErrCheckpointsAltBeforeCheckpoint, cp.AlternativeBlockAllowed(11, 9))
	assert.Nil(t, cp.AlternativeBlockAllowed(11, 10))
	assert.Nil(t, cp.AlternativeBlockAllowed(11, 11))
}
