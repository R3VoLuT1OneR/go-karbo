package p2p

import (
	"context"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"math/rand"
	"net"
)

type HostConfig struct {
	PeerId uint64
	BindAddr *net.TCPAddr
	Network *config.Network
}

type Host struct {
	Config HostConfig

	Dialer *net.Dialer

	PeerStore *PeerStore

	Context context.Context
}

func NewHost(cfg HostConfig) Host {
	h := Host{
		Config: cfg,
	}
	h.defaults()
	return h
}

func (h *Host) Start(ctx context.Context) error {
	h.Context = ctx

	for _, seed := range h.Config.Network.SeedNodes {

		peer, err := h.handshakeWithAddr(seed)
		if err != nil {
			return err
		}

		cmd, err := peer.Protocol.ReadCommand()
		if err != nil {
			return err
		}

		var notification NotificationTxPool
		if err := binary.Unmarshal(cmd.Payload, &notification); err != nil {
			return err
		}
		fmt.Println("notification", notification)

		pong, err := peer.Ping(h)
		if err != nil {
			return err
		}

		fmt.Println("pong", pong)
	}

	return nil
}

func (h *Host) handshakeWithAddr(addr string) (*Peer, error) {
	conn, err := h.Dialer.DialContext(h.Context, "tcp4", addr)
	if err != nil {
		return nil, err
	}

	peer := Peer{
		Protocol: &LevinProtocol{conn},
	}

	_, err = peer.Handshake(h)
	if err != nil {
		return nil, err
	}

	// TODO: verify that response with handle shake is proper

	return &peer, nil
}

func (h *Host) defaults() {
	if h.Config.PeerId == 0 {
		h.Config.PeerId = rand.Uint64()
	}

	if h.Dialer == nil {
		h.Dialer = &net.Dialer{
			// LocalAddr: h.Config.BindAddr,
		}
	}
}
