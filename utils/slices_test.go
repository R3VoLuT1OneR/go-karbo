package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceIsUniqueUint32(t *testing.T) {
	assert.True(t, SliceIsUniqueUint32(&[]uint32{1, 2, 3}))
	assert.False(t, SliceIsUniqueUint32(&[]uint32{1, 1, 1}))
}
