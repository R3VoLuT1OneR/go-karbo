package cryptonote

import (
	"encoding/hex"
	"github.com/r3volut1oner/go-karbo/crypto"
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
		assert.Equal(t, td.publicSpendKey, hex.EncodeToString(address.SpendPublicKey[:]))
		assert.Equal(t, td.publicViewKey, hex.EncodeToString(address.ViewPublicKey[:]))
		assert.Equal(t, uint64(111), address.Tag)
	}
}

func TestAddressGenerate(t *testing.T) {
	for _, td := range testAddressCompilation {

		publicSpendKeyBytes, _ := hex.DecodeString(td.publicSpendKey)
		publicViewKeyBytes, _ := hex.DecodeString(td.publicViewKey)

		var publicSpendKey crypto.PublicKey
		copy(publicSpendKey[:], publicSpendKeyBytes[:])

		var publicViewKey crypto.PublicKey
		copy(publicViewKey[:], publicViewKeyBytes[:])

		address := Generate(uint64(111), publicSpendKey, publicViewKey)

		assert.Equal(t, td.address, address.Base58)
		assert.Equal(t, td.publicSpendKey, hex.EncodeToString(address.SpendPublicKey[:]))
		assert.Equal(t, td.publicViewKey, hex.EncodeToString(address.ViewPublicKey[:]))
		assert.Equal(t, uint64(111), address.Tag)
	}
}
