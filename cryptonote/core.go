package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
)

type Core struct {
	Network *config.Network

	// TODO: Must be used DB instance for fetching the blocks.
	blocks map[Hash]*Block
}

var ErrBlockExists = errors.New("block already exists")

func NewCore(network *config.Network) (*Core, error) {
	// TODO: Load blocks from db?
	genesisBlock, err := GenerateGenesisBlock(network)
	if err != nil {
		return nil, err
	}

	core := &Core{
		Network: network,
		blocks: map[Hash]*Block{},
	}

	if err := core.AddBlock(genesisBlock); err != nil {
		return nil, err
	}

	return core, nil
}

func (c *Core) AddBlock(b *Block) error {
	h, err := b.Hash()

	if err != nil {
		return err
	}

	if c.HasBlock(h) {
		return ErrBlockExists
	}

	c.blocks[*h] = b

	return nil
}

func (c *Core) HasBlock(h *Hash) bool {
	_, ok := c.blocks[*h]

	return ok
}
