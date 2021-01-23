package cryptonote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testingData = []*struct {
	spendKey       string
	viewKey        string
	publicSpendKey string
	publicViewKey  string
}{
	{
		"6390482f5b3a1fe7fef34577b2cd0d14f12c075578e21ecaf48d1fbc300cf80b",
		"2d92d42406f972c51bce29af0d7ece15284c3decc8d15afa9d72ac76e0d07508",
		"711e2156025f8b8d66aeb2908e21a08f971a3b3b722de0e0876b68bcf0c71b74",
		"ba8e26760a9262408f4cf67cf0b5f4c3e69a8a07367b77149dac04834b300f29",
	},
}

func TestGenerateViewFromSpend(t *testing.T) {
	for _, td := range testingData {
		spendKey, _ := KeyFromHex(td.spendKey)
		viewKey, _ := KeyFromHex(td.viewKey)

		assert.Equal(t, viewKey, ViewFromSpend(spendKey))
	}
}

func TestGetPublicKey(t *testing.T) {
	for _, td := range testingData {
		spendKey, _ := KeyFromHex(td.spendKey)
		viewKey, _ := KeyFromHex(td.viewKey)

		assert.Equal(t, td.publicSpendKey, PublicFromPrivate(spendKey).Hex())
		assert.Equal(t, td.publicViewKey, PublicFromPrivate(viewKey).Hex())
	}
}
