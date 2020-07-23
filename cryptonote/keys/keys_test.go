package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateViewFromSpend(t *testing.T) {
	spendKey, err := PrivateFromHex("6390482f5b3a1fe7fef34577b2cd0d14f12c075578e21ecaf48d1fbc300cf80b")
	assert.Nil(t, err)

	viewKey, err := PrivateFromHex("2d92d42406f972c51bce29af0d7ece15284c3decc8d15afa9d72ac76e0d07508")
	assert.Nil(t, err)

	assert.Equal(t, viewKey, GenerateViewFromSpend(&spendKey))
}
