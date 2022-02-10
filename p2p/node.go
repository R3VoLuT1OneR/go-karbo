package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"io"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type HostConfig struct {
	PeerID   uint64
	BindAddr string
	Network  *config.Network

	ListenConfig *net.ListenConfig
}

type Node struct {
	Config     HostConfig
	Blockchain *cryptonote.BlockChain

	dialer *net.Dialer
	logger *log.Logger
	wg     *sync.WaitGroup
	ps     *peerStore

	context context.Context

	listener *net.TCPListener
}

// NewNode creates instance of the node
func NewNode(core *cryptonote.BlockChain, cfg HostConfig, logger *log.Logger) Node {
	var wg sync.WaitGroup

	h := Node{
		Config:     cfg,
		Blockchain: core,
		logger:     logger,
	}

	h.defaults()
	h.ps = NewPeerStore()
	h.wg = &wg

	return h
}

// Run the node server
// Listen for new connections on the defined port.
// Send handshake requests to the seed nodes.
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

	n.context = ctx

	n.wg.Add(1)
	go n.runListener()

	// TODO: Disabled only for check the listener
	//for _, seedAddr := range n.Config.Network.SeedNodes {
	//	go n.syncWithAddr(seedAddr)
	//}

	n.wg.Wait()
	return nil
}

// runListener listens for a new connections from external peers.
func (n *Node) runListener() {
	for {
		select {
		case <-n.context.Done():
			if err := n.listener.Close(); err != nil {
				n.logger.Errorf("failed to close listener: %s", err)
			}

			n.wg.Done()
			return
		default:
			_ = n.listener.SetDeadline(time.Now().Add(time.Second * 5))

			conn, err := n.listener.AcceptTCP()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					time.Sleep(time.Second)
					continue
				}
				n.logger.Errorf("failed to accept connection: %s", err)
			}

			go n.handleIncomingConnection(conn)
		}
	}
}

func (n *Node) handleIncomingConnection(conn *net.TCPConn) {
	// TODO: Enabling handling incoming connections
	//return

	peer := NewPeerFromIncomingConnection(conn)

	//// TODO: Add peer to peerstore. Make sure it is not exists.
	//
	n.wg.Add(1)
	defer n.wg.Done()

	n.listenForCommands(peer)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}

func (n *Node) listenForCommands(p *Peer) {
	for {
		// Peer state changes asynchronously after handling some commands.
		// Here we are taking care of handle different peer statuses.
		switch p.state {
		// Our node must be synchronized with the peer
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing

			if err := p.requestChain(n.Blockchain); err != nil {
				n.logger.Errorf("failed to write request chain: %s", err)
			}

		// Peer shutdown.
		// Stop listening for commands.
		case PeerStateShutdown:
			n.logger.Infof("[%d] shutting down...", p.ID)
			return
		}

		// Wait 1 second before read next command from the connection
		// Stop listening for the commands by return from the function
		select {
		case <-time.After(time.Second):
		case <-n.context.Done():
			return
		}

		// Read a next command from the connection
		cmd, err := p.protocol.read()

		// No command received continue to next loop cycle
		if err == io.EOF {
			continue
		}

		// On any error we move the peer to the grey list
		if err != nil {
			log.Errorf("error on read command: %s", err)
			_ = n.ps.toGrey(p)
			break
		}

		// There is special command type as "notification" we handle them with a separate method.
		if cmd.IsNotify {
			if err := n.handleNotification(p, cmd); err != nil {
				n.logger.Errorf("failed to handle notification %d: %s", cmd.Command, err)
			}

			continue
		}

		// Call method for handling the notification
		if err := n.handleCommand(p, cmd); err != nil {
			n.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
		}
	}
}

// handleNotification
//
// Receive notification from remote peer and handle it depend on the notification code.
func (n *Node) handleNotification(p *Peer, cmd *LevinCommand) error {
	n.logger.Tracef("[%s] handeling notification: %d", p, cmd.Command)

	// TODO: This part of the code can be used for debug
	//cwd, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//err = ioutil.WriteFile(
	//	fmt.Sprintf("%s/%d.dat", cwd, cmd.Command),
	//	cmd.Payload,
	//	0644,
	//)
	//if err != nil {
	//	panic(err)
	//}

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

		genesisBlock, err := n.Blockchain.GenesisBlock()
		if err != nil {
			return fmt.Errorf("[%s] unexpected error: %w", p, err)
		}

		// Make sure genesis blocks belongs to same network
		if *genesisBlock.Hash() != notification.Blocks[len(notification.Blocks)-1] {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] request chain genesis block not match", p))
		}

		// Get max index of this node blockchain
		//topIndex, err := p.node.Blockchain.TopIndex()
		//if err != nil {
		//	return fmt.Errorf("[%s] unexpected error: %w", p, err)
		//}

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
		hasFirstBlock := n.Blockchain.HaveBlock(&firstHash)

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
			hasBlock := n.Blockchain.HaveBlock(&bh)

			if allBlockKnown && hasBlock {
				continue
			}

			allBlockKnown = false
			p.neededBlocks = append(p.neededBlocks, bh)
		}

		return p.requestMissingBlocks(n.Blockchain, false)
	case NotificationResponseGetObjects:
		notification := nt.(NotificationResponseGetObjects)

		n.logger.Debugf(
			"[%s] NotificationResponseGetObjects, height: %d",
			p, notification.CurrentBlockchainHeight,
		)

		return p.handleResponseGetObjects(n.Blockchain, notification)
	default:
		n.logger.Errorf("can't handle notification type: %s", reflect.TypeOf(nt))
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
		if err := HandleHandshake(n, p, c.(HandshakeRequest), cmd); err != nil {
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

func (n *Node) syncWithAddr(addr string) {
	ctx, cancel := context.WithCancel(n.context)
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
		n.logger.Errorf("failed handshake: %s", err)
		cancel()
		return
	}

	n.logger.Debugf("[%s] handshake established", peer)

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

	n.listenForCommands(peer)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}
