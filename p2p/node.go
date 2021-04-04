package p2p

import (
	"context"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"math/rand"
	"net"
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

// NewNode creates instance of the node
func NewNode(core *cryptonote.Core, cfg HostConfig, logger *log.Logger) Node {
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

	peer.listenForCommands(ctx)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}
