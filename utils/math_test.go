package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMedianSlice(t *testing.T) {
	assert.Equal(t, uint64(0), MedianSlice([]uint64{}))
	assert.Equal(t, uint64(75), MedianSlice([]uint64{75}))
	assert.Equal(t, uint64(70), MedianSlice([]uint64{75, 66}))
	assert.Equal(t, uint64(5), MedianSlice([]uint64{1, 3, 5, 7, 9}))
}
