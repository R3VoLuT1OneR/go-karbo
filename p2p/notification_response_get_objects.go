package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"io/ioutil"
)

// NotificationResponseGetObjects == 2004
type NotificationResponseGetObjects struct {
	// Exists in old legacy code in definition but not exists in notification
	Transactions            []string      `binary:"txs,omitempty"`
	Blocks                  []RawBlock    `binary:"blocks,array"`
	CurrentBlockchainHeight uint32        `binary:"current_blockchain_height"`
	MissedIds               []crypto.Hash `binary:"missed_ids,binary,omitempty"`
}

func (n *Node) HandleResponseGetObjects(p *Peer, nt NotificationResponseGetObjects) error {
	p.logger.Debugf(
		"received response to get objects, height: %d",
		nt.CurrentBlockchainHeight,
	)

	if len(nt.Blocks) == 0 {
		p.Shutdown()
		return errors.New(fmt.Sprintf("[%s] got zero blocks on get objects", p))
	}

	if p.lastResponseHeight > nt.CurrentBlockchainHeight {
		p.Shutdown()
		return errors.New(fmt.Sprintf(
			"[%s] got wrong currentBlockchainHeight = %d, current = %d", p,
			nt.CurrentBlockchainHeight,
			p.lastResponseHeight,
		))
	}

	// TODO: Implement P2P Node observable height (max observed height) and update it if new observed height is found.

	p.remoteHeight = nt.CurrentBlockchainHeight

	orderedBlocks := make([]*cryptonote.Block, len(nt.Blocks))
	transactions := map[crypto.Hash][][]byte{}
	for i, rawBlock := range nt.Blocks {
		block := cryptonote.Block{}
		rawBlockReader := bytes.NewReader(rawBlock.Block)
		if err := block.Deserialize(rawBlockReader); err != nil {
			p.Shutdown()
			// TODO: Remove this. It is debug only.
			height := n.Blockchain.Height()
			blockHeight := height + uint32(i)
			_ = ioutil.WriteFile(fmt.Sprintf("./block_%d.dat", blockHeight), rawBlock.Block, 0644)
			return errors.New(
				fmt.Sprintf("[%s] (%d) failed to convert raw block (%d): %s", p, i, blockHeight, err),
			)
		}

		hash := block.Hash()
		if !p.requestedBlocks.Has(hash) {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] got not requested block #%d '%s'", p, i, hash.String()))
		}

		if len(block.TransactionsHashes) != len(rawBlock.Transactions) {
			p.Shutdown()
			return errors.New(fmt.Sprintf(
				"[%s] got wrong block transactions size. block: %s block tx: %d raw tx: %d",
				p, hash.String(), len(block.TransactionsHashes), len(rawBlock.Transactions),
			))
		}

		p.requestedBlocks.Remove(hash)
		transactions[*hash] = rawBlock.Transactions
		orderedBlocks[i] = &block
	}

	if len(p.requestedBlocks) > 0 {
		p.Shutdown()
		return errors.New(fmt.Sprintf(
			"[%s] got not all requested objectes, missing %d", p, len(p.requestedBlocks),
		))
	}

	if err := n.processNewObjects(orderedBlocks, transactions); err != nil {
		return err
	}

	height := n.Blockchain.Height()
	p.logger.Infof("process block, total height: %d", height)

	return p.requestMissingBlocks(n, true)
}
