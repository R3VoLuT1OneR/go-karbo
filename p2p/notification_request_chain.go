package p2p

import (
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
)

type NotificationRequestChain struct {
	Blocks []crypto.Hash `binary:"block_ids,binary"`
}

func NewRequestChain(bc *cryptonote.BlockChain) (*NotificationRequestChain, error) {
	list, err := bc.BuildSparseChain()
	if err != nil {
		return nil, err
	}

	return &NotificationRequestChain{list}, nil
}

func (n *Node) NotifyRequestChain(p *Peer) error {
	notification, err := NewRequestChain(n.Blockchain)
	if err != nil {
		return err
	}

	p.logger.Debugf("request chain %d (%d blocks) ", p.lastResponseHeight, len(notification.Blocks))

	if err := p.protocol.Notify(NotificationRequestChainID, *notification); err != nil {
		return err
	}

	return nil
}
