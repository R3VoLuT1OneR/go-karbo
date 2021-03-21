package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
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

	node    *Node
	version byte
	state   byte

	address *net.TCPAddr

	protocol *LevinProtocol

	remoteHeight uint32
	lastResponseHeight uint32

	neededBlocks 	cryptonote.HashList
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
		ID:       rand.Uint64(),
		node:     n,
		protocol: &LevinProtocol{conn},
		address:  tcpAddr,
	}

	return &peer, nil
}

func NewPeerFromIncomingConnection(conn net.Conn) *Peer {
	return &Peer{
		protocol: &LevinProtocol{conn},
	}
}

func (p *Peer) listenForCommands(ctx context.Context) {
	for {
		switch p.state {
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing
			if err := p.requestChain(p.node); err != nil {
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

		n.logger.Infof("txs pool: %d", len(notification.Transactions))
	case NotificationResponseChainEntry:
		notification := nt.(NotificationResponseChainEntry)

		n.logger.Tracef(
			"response chain entry: %d -> %d (%d)",
			notification.Start,
			notification.Total,
			len(notification.BlocksHashes),
		)

		if len(notification.BlocksHashes) == 0 {
			p.state = PeerStateShutdown
			return errors.New(fmt.Sprintf("[%d] received empty blocks in response chain enrty", p.ID))
		}

		// TODO: Assert first block is known to our blockchain

		p.remoteHeight = notification.Total
		p.lastResponseHeight = notification.Start + uint32(len(notification.BlocksHashes) - 1)

		if p.lastResponseHeight > p.remoteHeight {
			p.state = PeerStateShutdown
			return errors.New(
				fmt.Sprintf(
					"[%s] sent wrong response chain entry, with Total = %d, Start = %d, blocks = %d", p,
					notification.Start,
					notification.Total,
					len(notification.BlocksHashes),
				),
			)
		}

		allBlockKnown := true
		for _, bh := range notification.BlocksHashes {
			if allBlockKnown && n.Core.HasBlock(&bh) {
				continue
			}

			allBlockKnown = false
			p.neededBlocks = append(p.neededBlocks, bh)
		}

		return p.requestMissingBlocks(false)
	case NotificationResponseGetObjects:
		return p.handleResponseGetObjects(nt.(NotificationResponseGetObjects))
	default:
		n.logger.Errorf("can't handle notificaiton type: %s", reflect.TypeOf(n))
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
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf("[%s] got zero blocks on get objects", p))
	}

	if p.lastResponseHeight > nt.CurrentBlockchainHeight {
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf(
			"[%s] got wrong currentBlockchainHeight = %d, current = %d", p,
			nt.CurrentBlockchainHeight,
			p.lastResponseHeight,
		))
	}

	// TODO: Update observedHeight

	p.remoteHeight = nt.CurrentBlockchainHeight

	var blocks []cryptonote.Block
	for i, rawBlock := range nt.Blocks {
		block, err := rawBlock.ToBlock()
		if err != nil {
			p.state = PeerStateShutdown
			return errors.New(fmt.Sprintf("[%s] failed to convert raw block to block: %s", p, err))
		}

		hash, err := block.Hash()
		if err != nil {
			return err
		}

		if !p.requestedBlocks.Has(hash) {
			p.state = PeerStateShutdown

			ioutil.WriteFile(fmt.Sprintf("./block_%d.dat", i), rawBlock.Block, 0644)
			for ti, tbytes := range rawBlock.Transactions {
				ioutil.WriteFile(fmt.Sprintf("./block_%d_trans_%d.dat", i, ti), tbytes, 0644)
			}

			return errors.New(fmt.Sprintf("[%s] got not requested block #%d '%s'", p, i, hash.String()))
		}

		p.requestedBlocks.Remove(hash)
		blocks = append(blocks, *block)
	}

	if len(p.requestedBlocks) > 0 {
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf(
			"[%s] got not all requested objectes, missing %d", p, len(p.requestedBlocks),
		))
	}

	if err := p.processNewBlocks(blocks); err != nil {
		return err
	}

	p.node.logger.Infof("process block, total height: %d", p.node.Core.Height())

	return p.requestMissingBlocks(true)
}

func (p *Peer) processNewBlocks(blocks []cryptonote.Block) error {
	core := p.node.Core

	for _, block := range blocks {
		if err := core.AddBlock(&block); err != nil {
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

	req, err := NewHandshakeRequest(h.Config.Network)
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

	p.state = PeerStateSynchronizing
	p.version = res.NodeData.Version
	p.remoteHeight = res.PayloadData.CurrentHeight

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

func (p *Peer) requestChain(h *Node) error {
	requestChain, err := newRequestChain(h.Config.Network)
	if err != nil {
		return err
	}

	if err := p.protocol.Notify(NotificationRequestChainID, *requestChain); err != nil {
		return err
	}

	return nil
}

func (p *Peer) processSyncData(data SyncData, initial bool) error {
	// TODO: Implement sync data
	//if p.state == PeerStateBeforeHandshake && !initial {
	//	return nil
	//}

	p.state = PeerStateSyncRequired

	return nil
}

// TODO: Maybe move to different place?
func (p *Peer) requestMissingBlocks(checkHavingBlocks bool) error {
	if len(p.neededBlocks) > 0 {
		n := NotificationRequestGetObjects{}

		for len(p.neededBlocks) > 0 && len(n.Blocks) < MaxBlockSynchronization {
			nb := p.neededBlocks[0]

			if p.node == nil {
				fmt.Println("p.core", p)
			}
			if !(checkHavingBlocks && p.node.Core.HasBlock(&nb)) {
				n.Blocks = append(n.Blocks, nb)
				p.requestedBlocks = append(p.requestedBlocks, nb)
			}

			p.neededBlocks = p.neededBlocks[1:]
		}

		if err := p.protocol.Notify(NotificationRequestGetObjectsID, n); err != nil {
			return err
		}
	} else if p.lastResponseHeight < p.remoteHeight {
		// TODO: Send request chain
	} else {
		// TODO Check this condition
		//if (!(context.m_last_response_height ==
		//	context.m_remote_blockchain_height - 1 &&
		//	!context.m_needed_objects.size() &&
		//	!context.m_requested_objects.size())) {
		//	logger(Logging::ERROR, Logging::BRIGHT_RED)
		//	<< "request_missing_blocks final condition failed!"
		//	<< "\r\nm_last_response_height=" << context.m_last_response_height
		//	<< "\r\nm_remote_blockchain_height=" << context.m_remote_blockchain_height
		//	<< "\r\nm_needed_objects.size()=" << context.m_needed_objects.size()
		//	<< "\r\nm_requested_objects.size()=" << context.m_requested_objects.size()
		//	<< "\r\non connection [" << context << "]";

		// TODO: Request missing transactions

		p.state = PeerStateNormal
		// h.logger.Infof("[%s] syncronized", p)
		// TODO: On connection synchronized
	}

	return nil
}

func (p *Peer) String() string {
	return fmt.Sprintf("%d", p.ID)
}