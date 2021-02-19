package p2p

import (
	"bytes"
	"context"
	"encoding/binary"
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

type Host struct {
	Config HostConfig
	Core   *cryptonote.Core

	dialer   *net.Dialer
	logger   *log.Logger
	wg       *sync.WaitGroup
	ps 		 *peerStore

	listener *net.TCPListener
}

func NewHost(core *cryptonote.Core, cfg HostConfig, logger *log.Logger) Host {
	var wg sync.WaitGroup

	h := Host{
		Config: cfg,
		Core: core,
		logger: logger,
	}

	h.defaults()
	h.ps = NewPeerStore()
	h.wg = &wg

	return h
}

func (h *Host) defaults() {
	if h.Config.PeerID == 0 {
		h.Config.PeerID = rand.Uint64()
	}

	if h.Config.ListenConfig == nil {
		h.Config.ListenConfig = &net.ListenConfig{}
	}

	if h.dialer == nil {
		h.dialer = &net.Dialer{
			//LocalAddr: h.Config.BindAddr,
			Timeout: time.Second,
		}
	}
}

func (h *Host) Run(ctx context.Context) error {
	// listener, err := h.Config.ListenConfig.Listen(ctx, "tcp", h.Config.BindAddr)
	addr, err := net.ResolveTCPAddr("tcp", h.Config.BindAddr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	h.listener = listener
	h.logger.Debugf("listening on %s", listener.Addr())

	h.wg.Add(1)
	go h.runListener(ctx)

	for _, seedAddr := range h.Config.Network.SeedNodes {
		go h.syncWithAddr(ctx, seedAddr)
	}

	h.wg.Wait()
	return nil
}

func (h *Host) runListener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if err := h.listener.Close(); err != nil {
				h.logger.Errorf("failed to close listener: %s", err)
			}

			h.wg.Done()
			return
		default:
			_ = h.listener.SetDeadline(time.Now().Add(time.Second * 5))

			conn, err := h.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					time.Sleep(time.Second)
					continue
				}
				h.logger.Errorf("failed to accept connection: %s", err)
			}

			go h.handleIncomingConnection(ctx, conn)
		}
	}
}

func (h *Host) handleIncomingConnection(ctx context.Context, conn net.Conn) {
	peer := NewPeerFromIncomingConnection(conn)

	h.wg.Add(1)
	defer h.wg.Done()

	h.listenForCommands(ctx, peer)
}

func (h *Host) syncWithAddr(c context.Context, addr string) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	peer, err := NewPeerFromTCPAddress(ctx, h, addr)
	if err != nil {
		// h.logger.Errorf("failed to dial to peer: %s", err)
		cancel()
		return
	}

	//handshake, err := peer.handshake(h)
	_, err = peer.handshake(h)
	if err != nil {
		h.logger.Error("failed handshake")
		cancel()
		return
	}

	h.logger.Debugf("[#%16x] handshake established", peer.ID)

	if err := h.ps.toWhite(peer); err != nil {
		h.logger.Error("failed to add peer to the store")
		cancel()
		return
	}

	h.wg.Add(1)
	defer h.wg.Done()

	//for _, pe := range handshake.Peers {
	//	go h.syncWithAddr(c, pe.Address.String())
	//}

	h.listenForCommands(ctx, peer)

	if err := h.ps.toGrey(peer); err != nil {
		h.logger.Warnf("peer remove failed: %s", err)
	}

	h.logger.Debugf("[%16x] sync closed", peer.ID)
}

func (h *Host) listenForCommands(ctx context.Context, p *Peer) {
	for {
		switch p.state {
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing
			if err := p.requestChain(h); err != nil {
				h.logger.Errorf("failed to write request chain: %s", err)
			}

		case PeerStateShutdown:
			h.logger.Infof("[%d] shutting down...", p.ID)
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
			_ = h.ps.toGrey(p)
			break
		}

		if cmd.IsNotify {
			if err := h.handleNotification(p, cmd); err != nil {
				h.logger.Errorf("failed to handle notification %d: %s", cmd.Command, err)
			}

			continue
		}

		if err := h.handleCommand(p, cmd); err != nil {
			h.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
		}
	}
}

func (h *Host) handleNotification(p *Peer, cmd *LevinCommand) error {
	h.logger.Tracef("[%s] handeling notification: %d", p, cmd.Command)

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

	n, err := parseNotification(cmd)
	if err != nil {
		return err
	}

	switch n.(type) {
	case NotificationTxPool:
		notification := n.(NotificationTxPool)

		h.logger.Infof("txs pool: %d", len(notification.Transactions))
	case NotificationResponseChainEntry:
		notification := n.(NotificationResponseChainEntry)

		h.logger.Tracef(
			"response chain entry: %d -> %d (%d)",
			notification.Start,
			notification.Total,
			len(notification.BlockIds),
		)

		if len(notification.BlockIds) == 0 {
			p.state = PeerStateShutdown
			return errors.New(fmt.Sprintf("[%d] received empty blocks in response chain enrty", p.ID))
		}

		// TODO: Assert first block is known to our blockchain

		p.remoteHeight = notification.Total
		p.lastResponseHeight = notification.Start + uint32(len(notification.BlockIds) - 1)

		if p.lastResponseHeight > p.remoteHeight {
			p.state = PeerStateShutdown
			return errors.New(
				fmt.Sprintf(
					"[%s] sent wrong response chain entry, with Total = %d, Start = %d, blocks = %d", p,
					notification.Start,
					notification.Total,
					len(notification.BlockIds),
				),
			)
		}

		allBlockKnown := true
		for _, bh := range notification.BlockIds {
			if allBlockKnown && h.Core.HasBlock(&bh) {
				continue
			}

			allBlockKnown = false
			p.neededBlocks = append(p.neededBlocks, bh)
		}

		return p.requestMissingBlocks(false)
	case NotificationResponseGetObjects:
		return h.handleResponseGetObjects(p, n.(NotificationResponseGetObjects))
	default:
		h.logger.Errorf("can't handle notificaiton type: %s", reflect.TypeOf(n))
	}

	return nil
}

func (h *Host) handleCommand(p *Peer, cmd *LevinCommand) error {
	c, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	switch c.(type) {
	case HandshakeRequest:
		// TODO: Check peer network and rest of the data
		handshakeRequest := c.(HandshakeRequest)
		if handshakeRequest.NodeData.NetworkID != h.Config.Network.NetworkID {
			return errors.New("wrong network on handshake")
		}

		// TODO: Send ping and make sure we can connect to the peer and add it to the white list.
		//if err := p.processSyncData(c.(HandshakeRequest).PayloadData, true); err != nil {
		//	return err
		//}

		h.logger.Debugf("[%v] handshake received", p.ID)

		rsp, err := NewHandshakeResponse(h)
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

		res, err := newTimedSyncResponse(h)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *res, 1); err != nil {
			return err
		}

		h.logger.Infof("[%s] sync request %d", p, command.PayloadData.CurrentHeight)
	default:
		h.logger.Errorf("received unknown commands type: %s", reflect.TypeOf(c))
	}

	return nil
}

func (h *Host) handleResponseGetObjects(p *Peer, n NotificationResponseGetObjects) error {

	h.logger.Tracef("[%s] response to get objects", p)

	if len(n.Blocks) == 0 {
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf("[%s] got zer blocks on get objects", p))
	}

	if p.lastResponseHeight > n.CurrentBlockchainHeight {
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf(
			"[%s] got wrong currentBlockchainHeight = %d, current = %d", p,
			n.CurrentBlockchainHeight,
			p.lastResponseHeight,
		))
	}

	// TODO: Update observedHeight

	p.remoteHeight = n.CurrentBlockchainHeight

	var blocks []cryptonote.Block
	for _, rawBlock := range n.Blocks {
		var block cryptonote.Block
		reader := bytes.NewReader(rawBlock.Block)
		if err := binary.Read(reader, binary.LittleEndian, block); err != nil {
			p.state = PeerStateShutdown
			return errors.New(fmt.Sprintf("[%s] faield to convert raw block to block", p))
		}

		// TODO: Set idle

		hash, err := block.Hash()
		if err != nil {
			return err
		}

		if !p.requestedBlocks.Has(hash) {
			p.state = PeerStateShutdown
			return errors.New(fmt.Sprintf("[%s] got not requested block '%s'", p, hash.String()))
		}

		p.requestedBlocks.Remove(hash)
		blocks = append(blocks, block)
	}

	if len(p.requestedBlocks) > 0 {
		p.state = PeerStateShutdown
		return errors.New(fmt.Sprintf(
			"[%s] got not all requested objectes, missing %d", p, len(p.requestedBlocks),
		))
	}

	if err := h.processBlocks(p, blocks); err != nil {
		return err
	}

	return p.requestMissingBlocks(true)
}

func (h *Host) processBlocks(p *Peer, blocks []cryptonote.Block) error {
	for _, block := range blocks {
		if err := h.Core.AddBlock(&block); err != nil {
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
