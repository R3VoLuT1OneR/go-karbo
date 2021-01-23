package cryptonote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testAddressCompilation = []*struct {
	address        string
	publicSpendKey string
	publicViewKey  string
}{
	{
		"KiAe1HTejHad5mgdQGSxuWZAFuCLo7KkhVJkHTwrxMFSWHgj1FvajVNXsAwjo7PYdpBon3qJREB7iMDAGWCtqvRjFoCrBVD",
		"dbbd38087aa091d7b8bf36d30d4059c0450ac73c9ac5dea93af77d1e6cd7fdaf",
		"197b9cfc60e48db887a94cc359c2c3409b998b88d7241578d35b41a66df15e83",
	},
}

func TestAddressFromString(t *testing.T) {
	for _, td := range testAddressCompilation {
		address, err := FromString(td.address)

		assert.Nil(t, err)
		assert.Equal(t, td.address, address.Base58)
		assert.Equal(t, td.publicSpendKey, address.SpendPublicKey.Hex())
		assert.Equal(t, td.publicViewKey, address.ViewPublicKey.Hex())
		assert.Equal(t, uint64(111), address.Tag)
	}
}

func TestAddressGenerate(t *testing.T) {
	for _, td := range testAddressCompilation {
		publicSpendKey, _ := KeyFromHex(td.publicSpendKey)
		publicViewKey, _ := KeyFromHex(td.publicViewKey)
		address := Generate(uint64(111), publicSpendKey, publicViewKey)

		assert.Equal(t, td.address, address.Base58)
		assert.Equal(t, td.publicSpendKey, address.SpendPublicKey.Hex())
		assert.Equal(t, td.publicViewKey, address.ViewPublicKey.Hex())
		assert.Equal(t, uint64(111), address.Tag)
	}
}
