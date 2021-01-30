package p2p

import (
	"context"
	"github.com/r3volut1oner/go-karbo/config"
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

	ListenConfig  *net.ListenConfig
}

type Host struct {
	Config HostConfig

	dialer   *net.Dialer
	logger   *log.Logger
	wg       *sync.WaitGroup
	ps 		 *PeerStore

	listener *net.TCPListener
}

func NewHost(cfg HostConfig, logger *log.Logger) Host {
	var wg sync.WaitGroup

	h := Host{
		Config: cfg,
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
	h.logger.Infof("listening on %s", listener.Addr())

	h.wg.Add(1)
	go h.startListen(ctx)
	go h.startPeerSync(ctx)

	h.wg.Wait()
	return nil
}

func (h *Host) startListen(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			if err := h.listener.Close(); err != nil {
				h.logger.Errorf("failed to close listener: %s", err)
			}

			h.wg.Done()
			return
		default:
			_ = h.listener.SetDeadline(time.Now().Add(time.Second))
			_, err := h.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					time.Sleep(time.Second)
					continue
				}
				h.logger.Errorf("failed to accept connection: %s", err)
			}

			// TODO: Implement new connections listener
		}
	}
}

func (h *Host) startPeerSync(ctx context.Context) {
	h.logger.Print("sync started")

	for _, seedAddr := range h.Config.Network.SeedNodes {
		go h.syncWithAddr(ctx, seedAddr)
	}
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

	handshake, err := peer.Handshake(h)
	if err != nil {
		h.logger.Error("failed handshake")
		cancel()
		return
	}

	h.logger.Infof("[#%16x] handshake established with", peer.ID)

	if err := h.ps.Add(peer); err != nil {
		h.logger.Error("failed to add peer to the store")
		cancel()
		return
	}

	h.wg.Add(1)
	defer h.wg.Done()

	for _, pe := range handshake.Peers {
		go h.syncWithAddr(c, pe.Address.String())
	}

	h.listenNotifications(ctx, peer)

	if err := h.ps.Remove(peer); err != nil {
		h.logger.Warnf("peer remove failed: %s", err)
	}

	h.logger.Infof("[%16x] sync closed", peer.ID)
}

func (h *Host) listenNotifications(ctx context.Context, p *Peer) {
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}

		cmd, err := p.protocol.ReadCommand()
		if err == io.EOF {
			continue
		}

		if err != nil {
			log.Errorf("error on read command: %s", err)
		}

		if cmd.IsNotify {
			if err := h.handleNotification(cmd); err != nil {
				h.logger.Errorf("failed handle notification: %s", err)
			}

			continue
		}

		if err := h.handleCommand(p, cmd); err != nil {
			h.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
		}
	}
}

func (h *Host) handleNotification(cmd *LevinCommand) error {
	n, err := parseNotification(cmd)
	if err != nil {
		return err
	}

	switch n.(type) {
	case NotificationTxPool:
		h.logger.Infof("txs pool: %d", len(n.(NotificationTxPool).Transactions))
	default:
		h.logger.Errorf("received unknown notificaiton type: %s", reflect.TypeOf(n))
	}
	return nil
}

func (h *Host) handleCommand(p *Peer, cmd *LevinCommand) error {
	c, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	switch c.(type) {
	case TimedSyncRequest:
		// TODO: Handle sync request
		request := c.(TimedSyncRequest)

		res, err := newTimedSyncResponse(h.Config.Network)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *res, 1); err != nil {
			return err
		}

		h.logger.Infof("sync request %d", request.PayloadData.CurrentHeight)
	default:
		h.logger.Errorf("received unknown commands type: %s", reflect.TypeOf(c))
	}

	return nil
}
