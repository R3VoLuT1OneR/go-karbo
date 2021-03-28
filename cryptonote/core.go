package cryptonote

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
)

type Core struct {
	Network *config.Network

	storage Store
}

func NewCore(network *config.Network, DB Store) (*Core, error) {
	core := &Core{
		Network: network,
		storage: DB,
	}

	if err := core.Init(); err != nil {
		return nil, err
	}

	return core, nil
}

func (c *Core) Init() error {
	if err := c.initDB(); err != nil {
		return errors.New(fmt.Sprintf("failed to initialize storage: %s", err))
	}

	return nil
}

func (c *Core) AddBlock(b *Block) error {
	if err := c.storage.AppendBlock(b); err != nil {
		return err
	}

	return nil
}

func (c *Core) HasBlock(h *Hash) (bool, error) {
	hasBlock, err := c.storage.HasBlock(h)
	if err != nil {
		return false, nil
	}

	return hasBlock, nil
}

func (c *Core) Height() (uint32, error) {
	height, err := c.storage.GetHeight()
	if err != nil {
		return 0, err
	}

	return height, err
}

func (c *Core) TopBlock() (*Block, uint32, error) {
	height, err := c.Height()
	if err != nil {
		return nil, 0, err
	}

	block, err := c.storage.GetBlockByHeight(height)
	if err != nil {
		return nil, 0, err
	}

	return block, height, nil
}

func (c *Core) BuildSparseChain() ([]Hash, error) {
	var list []Hash

	topBlock, height, err := c.TopBlock()
	if err != nil {
		return nil, err
	}

	topHash, err := topBlock.Hash()
	if err != nil {
		return nil, err
	}

	list = append(list, *topHash)

	for i := uint32(1); i < height; i *= 2 {
		hash, err := c.storage.GetBlockHashByHeight(i)
		if err != nil {
			return nil, err
		}

		list = append(list, *hash)
	}

	ghash, err := c.storage.GetBlockHashByHeight(0)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(ghash[:], list[0][:]) && !bytes.Equal(ghash[:], list[len(list) - 1][:]) {
		list = append(list, *ghash)
	}

	return list, nil
}

func (c *Core) initDB() error {
	if err := c.storage.Init(c.Network); err != nil {
		return err
	}

	return nil
}