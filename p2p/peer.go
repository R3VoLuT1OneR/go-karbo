package p2p

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"time"
)

const (
	PeerStateBeforeHandshake byte = iota
	PeerStateSynchronizing
	PeerStateIdle
	PeerStateNormal
	PeerStateSyncRequired
	PeerStatePoolSyncRequired
	PeerStateShutdown
)

type Peer struct {
	// TODO: Make sure we generate local peer ID and updating external IDs
	ID uint64

	// Node to which peer connected to. Our main server node.
	node    *Node
	version byte
	state   byte

	address *net.TCPAddr

	protocol *LevinProtocol

	remoteHeight       uint32
	lastResponseHeight uint32

	neededBlocks    cryptonote.HashList
	requestedBlocks cryptonote.HashList
}

func NewPeerFromTCPAddress(ctx context.Context, n *Node, addr string) (*Peer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := n.dialer.DialContext(ctx, "tcp4", addr)
	if err != nil {
		return nil, err
	}

	peer := Peer{
		node:     n,
		protocol: &LevinProtocol{conn},
		address:  tcpAddr,
	}

	return &peer, nil
}

func NewPeerFromIncomingConnection(node *Node, conn net.Conn) *Peer {
	return &Peer{
		node:     node,
		protocol: &LevinProtocol{conn},
	}
}

func (p *Peer) Shutdown() {
	p.state = PeerStateShutdown
}

func (p *Peer) String() string {
	return fmt.Sprintf("%d", p.ID)
}

func (p *Peer) listenForCommands(ctx context.Context) {
	for {
		switch p.state {
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing

			if err := p.requestChain(); err != nil {
				p.node.logger.Errorf("failed to write request chain: %s", err)
			}

		case PeerStateShutdown:
			p.node.logger.Infof("[%d] shutting down...", p.ID)
			return
		}

		select {
		case <-time.After(time.Second * 3):
		case <-ctx.Done():
			return
		}

		cmd, err := p.protocol.read()
		if err == io.EOF {
			continue
		}

		if err != nil {
			log.Errorf("error on read command: %s", err)
			_ = p.node.ps.toGrey(p)
			break
		}

		if cmd.IsNotify {
			if err := p.node.handleNotification(p, cmd); err != nil {
				p.node.logger.Errorf("failed to handle notification %d: %s", cmd.Command, err)
			}

			continue
		}

		if err := p.handleCommand(cmd); err != nil {
			p.node.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
		}
	}
}

// handleNotification
//
// Receive notification from remote peer and handle it depend on the notification code.
func (n *Node) handleNotification(p *Peer, cmd *LevinCommand) error {
	n.logger.Tracef("[%s] handeling notification: %d", p, cmd.Command)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(
		fmt.Sprintf("%s/%d.dat", cwd, cmd.Command),
		cmd.Payload,
		0644,
	)
	if err != nil {
		panic(err)
	}

	nt, err := parseNotification(cmd)
	if err != nil {
		return err
	}

	switch nt.(type) {
	case NotificationTxPool:
		notification := nt.(NotificationTxPool)

		n.logger.Debugf("[%s] notification tx pool, size: %d", p, len(notification.Transactions))
	case NotificationRequestChain:
		notification := nt.(NotificationRequestChain)

		// Verify more than 0 requested blocks
		if len(notification.Blocks) == 0 {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] request chain with 0 blocks", p))
		}

		genesisBlockHash, err := p.node.Core.GenesisBlockHash()
		if err != nil {
			return fmt.Errorf("[%s] unexpected error: %w", p, err)
		}

		// Make sure genesis blocks belongs to same network
		if *genesisBlockHash != notification.Blocks[len(notification.Blocks)-1] {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] request chain genesis block not match", p))
		}

		// Get max index of this node blockchain
		topIndex, err := p.node.Core.TopIndex()
		if err != nil {
			return fmt.Errorf("[%s] unexpected error: %w", p, err)
		}

		// TODO: Build response chain entry and response to requested peer
		//responseChainEntry := &NotificationResponseChainEntry{
		//	TotalHeight: topIndex + 1,
		//}

		n.logger.Tracef("[%s] request chain %d blocks.", p, len(notification.Blocks))
	case NotificationResponseChainEntry:
		notification := nt.(NotificationResponseChainEntry)

		n.logger.Tracef(
			"[%s] notification response chain entry, start: %d, total: %d, blocks: %d",
			p, notification.StartHeight, notification.TotalHeight, len(notification.BlocksHashes),
		)

		if len(notification.BlocksHashes) == 0 {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] received empty blocks in response chain enrty", p))
		}

		firstHash := notification.BlocksHashes[0]
		hasFirstBlock, err := n.Core.HasBlock(&firstHash)
		if err != nil {
			return err
		}
		if !hasFirstBlock {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] hash %s missing in our blockchain", p, firstHash.String()))
		}

		p.remoteHeight = notification.TotalHeight
		p.lastResponseHeight = notification.StartHeight + uint32(len(notification.BlocksHashes)-1)

		if p.lastResponseHeight > p.remoteHeight {
			p.Shutdown()
			return errors.New(
				fmt.Sprintf(
					"[%s] sent wrong response chain entry, with TotalHeight = %d, StartHeight = %d, blocks = %d", p,
					notification.StartHeight,
					notification.TotalHeight,
					len(notification.BlocksHashes),
				),
			)
		}

		allBlockKnown := true
		for _, bh := range notification.BlocksHashes {
			hasBlock, err := n.Core.HasBlock(&bh)
			if err != nil {
				return err
			}

			if allBlockKnown && hasBlock {
				continue
			}

			allBlockKnown = false
			p.neededBlocks = append(p.neededBlocks, bh)
		}

		return p.requestMissingBlocks(false)
	case NotificationResponseGetObjects:
		notification := nt.(NotificationResponseGetObjects)

		n.logger.Debugf(
			"[%s] NotificationResponseGetObjects, height: %d",
			p, notification.CurrentBlockchainHeight,
		)

		return p.handleResponseGetObjects(notification)
	default:
		n.logger.Errorf("can't handle notification type: %s", reflect.TypeOf(nt))
	}

	return nil
}

func (p *Peer) handleCommand(cmd *LevinCommand) error {
	c, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	switch c.(type) {
	case HandshakeRequest:
		// TODO: Check peer network and rest of the data
		handshakeRequest := c.(HandshakeRequest)
		if handshakeRequest.NodeData.NetworkID != p.node.Config.Network.NetworkID {
			return errors.New("wrong network on handshake")
		}

		// TODO: Send ping and make sure we can connect to the peer and add it to the white list.
		//if err := p.processSyncData(c.(HandshakeRequest).PayloadData, true); err != nil {
		//	return err
		//}

		p.node.logger.Debugf("[%v] handshake received", p.ID)

		rsp, err := NewHandshakeResponse(p.node)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *rsp, 1); err != nil {
			return err
		}

	case TimedSyncRequest:
		command := c.(TimedSyncRequest)
		if err := p.processSyncData(command.PayloadData, false); err != nil {
			return err
		}

		res, err := newTimedSyncResponse(p.node)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *res, 1); err != nil {
			return err
		}

		p.node.logger.Infof("[%s] sync request %d", p, command.PayloadData.CurrentHeight)
	default:
		p.node.logger.Errorf("received unknown commands type: %s", reflect.TypeOf(c))
	}

	return nil
}

func (p *Peer) handleResponseGetObjects(nt NotificationResponseGetObjects) error {

	p.node.logger.Tracef("[%s] response to get objects", p)

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

	newObjects := map[*cryptonote.Block][][]byte{}
	for i, rawBlock := range nt.Blocks {
		block := cryptonote.Block{}
		rawBlockReader := bytes.NewReader(rawBlock.Block)
		if err := block.Deserialize(rawBlockReader); err != nil {
			p.Shutdown()
			height, _ := p.node.Core.TopIndex()
			blockHeight := int(height) + 1 + i
			_ = ioutil.WriteFile(fmt.Sprintf("./block_%d.dat", blockHeight), rawBlock.Block, 0644)
			return errors.New(
				fmt.Sprintf("[%s] (%d) failed to convert raw block (%d): %s", p, i, blockHeight, err),
			)
		}

		hash, err := block.Hash()
		if err != nil {
			return err
		}

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
		newObjects[&block] = rawBlock.Transactions
	}

	if len(p.requestedBlocks) > 0 {
		p.Shutdown()
		return errors.New(fmt.Sprintf(
			"[%s] got not all requested objectes, missing %d", p, len(p.requestedBlocks),
		))
	}

	if err := p.processNewObjects(newObjects); err != nil {
		return err
	}

	height, err := p.node.Core.TopIndex()
	if err != nil {
		return err
	}

	p.node.logger.Infof("process block, total height: %d", height)

	return p.requestMissingBlocks(true)
}

func (p *Peer) processNewObjects(objects map[*cryptonote.Block][][]byte) error {
	core := p.node.Core

	for block, transactions := range objects {
		if err := core.AddBlock(block, transactions); err != nil {
			return err
			// TODO: Process proper error
			//
			//if (addResult == error::AddBlockErrorCondition::BLOCK_VALIDATION_FAILED ||
			//	addResult == error::AddBlockErrorCondition::TRANSACTION_VALIDATION_FAILED ||
			//	addResult == error::AddBlockErrorCondition::DESERIALIZATION_FAILED) {
			//	logger(Logging::DEBUGGING) << context << "Block verification failed, dropping connection: " << addResult.message();
			//	m_p2p->drop_connection(context, true);
			//	return 1;
			//} else if (addResult == error::AddBlockErrorCondition::BLOCK_REJECTED) {
			//	logger(Logging::DEBUGGING) << context << "Block received at sync phase was marked as orphaned, dropping connection: " << addResult.message();
			//	m_p2p->drop_connection(context, true);
			//	return 1;
			//} else if (addResult == error::AddBlockErrorCode::ALREADY_EXISTS) {
			//	logger(Logging::DEBUGGING) << context << "Block already exists, switching to idle state: " << addResult.message();
			//	context.m_state = CryptoNoteConnectionContext::state_idle;
			//	context.m_needed_objects.clear();
			//	context.m_requested_objects.clear();
			//	return 1;
			//}
		}
	}

	return nil
}

func (p *Peer) handshake(h *Node) (*HandshakeResponse, error) {
	if p.state != PeerStateBeforeHandshake {
		return nil, errors.New("state is not before handshake")
	}

	req, err := NewHandshakeRequest(h.Core)
	if err != nil {
		return nil, err
	}

	var res HandshakeResponse
	if err := p.protocol.Invoke(CommandHandshake, *req, &res); err != nil {
		return nil, err
	}

	if h.Config.Network.NetworkID != res.NodeData.NetworkID {
		return nil, errors.New("wrong network id received")
	}

	if h.Config.Network.P2PMinimumVersion < res.NodeData.Version {
		return nil, errors.New("node data version not match minimal")
	}

	if err := p.processSyncData(res.PayloadData, true); err != nil {
		return nil, err
	}

	p.version = res.NodeData.Version
	p.ID = res.NodeData.PeerID

	// TODO: Handle new peerlist

	return &res, nil
}

func (p *Peer) ping() (*PingResponse, error) {
	req := PingRequest{}
	res := PingResponse{}

	if err := p.protocol.Invoke(CommandPing, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *Peer) requestChain() error {
	n, err := newRequestChain(p.node.Core)
	if err != nil {
		return err
	}

	p.node.logger.Tracef("[%s] request chain %d (%d blocks) ", p, p.lastResponseHeight, len(n.Blocks))

	if err := p.protocol.Notify(NotificationRequestChainID, *n); err != nil {
		return err
	}

	return nil
}

func (p *Peer) processSyncData(data SyncData, initial bool) error {
	// TODO: Implement sync data
	if p.state == PeerStateBeforeHandshake && !initial {
		return nil
	}

	if p.state == PeerStateSynchronizing {
	} else {
		p.state = PeerStateSyncRequired
	}

	p.remoteHeight = data.CurrentHeight

	return nil
}

func (p *Peer) requestMissingBlocks(checkHavingBlocks bool) error {
	if len(p.neededBlocks) > 0 {
		neededBlocks := p.neededBlocks
		requestBlocks := cryptonote.HashList{}

		for len(neededBlocks) > 0 && len(requestBlocks) < MaxBlockSynchronization {
			nb := neededBlocks[0]

			hasBlock, err := p.node.Core.HasBlock(&nb)
			if err != nil {
				return err
			}

			if !(checkHavingBlocks && hasBlock) {
				requestBlocks = append(requestBlocks, nb)
			}

			neededBlocks = neededBlocks[1:]
		}

		if len(requestBlocks) > 0 {
			n := NotificationRequestGetObjects{}
			n.Blocks = requestBlocks

			if err := p.protocol.Notify(NotificationRequestGetObjectsID, n); err != nil {
				return err
			}

			p.requestedBlocks = append(p.requestedBlocks, requestBlocks...)
		}

		p.neededBlocks = neededBlocks
	} else if p.lastResponseHeight < (p.remoteHeight - 1) {
		if err := p.requestChain(); err != nil {
			return err
		}
	} else {
		if p.lastResponseHeight == p.remoteHeight-1 && len(p.neededBlocks) != 0 && len(p.requestedBlocks) != 0 {
			return errors.New(fmt.Sprintf(
				"request missing blocks final condition failed: \n"+
					"response height: %d\n"+
					"remote blockchain height: %d\n"+
					"needed objects size: %d\n"+
					"requested objects size: %d",
				p.lastResponseHeight,
				p.remoteHeight,
				len(p.neededBlocks),
				len(p.requestedBlocks),
			))
		}

		// TODO: Request missing pool transactions
		// src/CryptoNoteProtocol/CryptoNoteProtocolHandler.cpp:907

		p.state = PeerStateNormal
		p.node.logger.Tracef("[%s] syncronized", p)

		// TODO: On connection synchronized
		// src/CryptoNoteProtocol/CryptoNoteProtocolHandler.cpp:911
	}

	return nil
}
