package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type HostConfig struct {
	PeerID   uint64
	BindAddr string
	Network  *config.Network

	ListenConfig  *net.ListenConfig
}

type Node struct {
	Config HostConfig
	Core   *cryptonote.Core

	dialer   *net.Dialer
	logger   *log.Logger
	wg       *sync.WaitGroup
	ps 		 *peerStore

	listener *net.TCPListener
}

func NewHost(core *cryptonote.Core, cfg HostConfig, logger *log.Logger) Node {
	var wg sync.WaitGroup

	h := Node{
		Config: cfg,
		Core: core,
		logger: logger,
	}

	h.defaults()
	h.ps = NewPeerStore()
	h.wg = &wg

	return h
}

func (n *Node) defaults() {
	if n.Config.PeerID == 0 {
		n.Config.PeerID = rand.Uint64()
	}

	if n.Config.ListenConfig == nil {
		n.Config.ListenConfig = &net.ListenConfig{}
	}

	if n.dialer == nil {
		n.dialer = &net.Dialer{
			//LocalAddr: n.Config.BindAddr,
			Timeout: time.Second,
		}
	}
}

func (n *Node) Run(ctx context.Context) error {
	// listener, err := n.Config.ListenConfig.Listen(ctx, "tcp", n.Config.BindAddr)
	addr, err := net.ResolveTCPAddr("tcp", n.Config.BindAddr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	n.listener = listener
	n.logger.Debugf("listening on %s", listener.Addr())

	n.wg.Add(1)
	go n.runListener(ctx)

	for _, seedAddr := range n.Config.Network.SeedNodes {
		go n.syncWithAddr(ctx, seedAddr)
	}

	n.wg.Wait()
	return nil
}

func (n *Node) runListener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if err := n.listener.Close(); err != nil {
				n.logger.Errorf("failed to close listener: %s", err)
			}

			n.wg.Done()
			return
		default:
			_ = n.listener.SetDeadline(time.Now().Add(time.Second * 5))

			conn, err := n.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					time.Sleep(time.Second)
					continue
				}
				n.logger.Errorf("failed to accept connection: %s", err)
			}

			go n.handleIncomingConnection(ctx, conn)
		}
	}
}

func (n *Node) handleIncomingConnection(ctx context.Context, conn net.Conn) {
	// TODO: Enabling handeling incomming connections
	return

	//peer := NewPeerFromIncomingConnection(conn)
	//
	//// TODO: Add peer to peerstore. Make sure it is not exists.
	//
	//n.wg.Add(1)
	//defer n.wg.Done()
	//
	//n.listenForCommands(ctx, peer)
}

func (n *Node) syncWithAddr(c context.Context, addr string) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	peer, err := NewPeerFromTCPAddress(ctx, n, addr)
	if err != nil {
		// n.logger.Errorf("failed to dial to peer: %s", err)
		cancel()
		return
	}

	//handshake, err := peer.handshake(n)
	_, err = peer.handshake(n)
	if err != nil {
		n.logger.Error("failed handshake")
		cancel()
		return
	}

	n.logger.Debugf("[#%16x] handshake established", peer.ID)

	if err := n.ps.toWhite(peer); err != nil {
		n.logger.Error("failed to add peer to the store")
		cancel()
		return
	}

	n.wg.Add(1)
	defer n.wg.Done()

	//for _, pe := range handshake.Peers {
	//	go n.syncWithAddr(c, pe.Address.String())
	//}

	n.listenForCommands(ctx, peer)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}

func (n *Node) listenForCommands(ctx context.Context, p *Peer) {
	for {
		switch p.state {
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing
			if err := p.requestChain(n); err != nil {
				n.logger.Errorf("failed to write request chain: %s", err)
			}

		case PeerStateShutdown:
			n.logger.Infof("[%d] shutting down...", p.ID)
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
			_ = n.ps.toGrey(p)
			break
		}

		if cmd.IsNotify {
			if err := n.handleNotification(p, cmd); err != nil {
				n.logger.Errorf("failed to handle notification %d: %s", cmd.Command, err)
			}

			continue
		}

		if err := n.handleCommand(p, cmd); err != nil {
			n.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
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
		return n.handleResponseGetObjects(p, nt.(NotificationResponseGetObjects))
	default:
		n.logger.Errorf("can't handle notificaiton type: %s", reflect.TypeOf(n))
	}

	return nil
}

func (n *Node) handleCommand(p *Peer, cmd *LevinCommand) error {
	c, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	switch c.(type) {
	case HandshakeRequest:
		// TODO: Check peer network and rest of the data
		handshakeRequest := c.(HandshakeRequest)
		if handshakeRequest.NodeData.NetworkID != n.Config.Network.NetworkID {
			return errors.New("wrong network on handshake")
		}

		// TODO: Send ping and make sure we can connect to the peer and add it to the white list.
		//if err := p.processSyncData(c.(HandshakeRequest).PayloadData, true); err != nil {
		//	return err
		//}

		n.logger.Debugf("[%v] handshake received", p.ID)

		rsp, err := NewHandshakeResponse(n)
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

		res, err := newTimedSyncResponse(n)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *res, 1); err != nil {
			return err
		}

		n.logger.Infof("[%s] sync request %d", p, command.PayloadData.CurrentHeight)
	default:
		n.logger.Errorf("received unknown commands type: %s", reflect.TypeOf(c))
	}

	return nil
}

func (n *Node) handleResponseGetObjects(p *Peer, nt NotificationResponseGetObjects) error {

	n.logger.Tracef("[%s] response to get objects", p)

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

	if err := n.processBlocks(blocks); err != nil {
		return err
	}

	n.logger.Infof("process block, total height: %d", n.Core.Height())

	return p.requestMissingBlocks(true)
}

func (n *Node) processBlocks(blocks []cryptonote.Block) error {
	for _, block := range blocks {
		if err := n.Core.AddBlock(&block); err != nil {
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
