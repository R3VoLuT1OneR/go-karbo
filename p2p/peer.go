package p2p

import (
	"context"
	"errors"
	"math/rand"
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
	ID uint64

	Version byte
	Height uint32
	State byte

	protocol *LevinProtocol
}

func NewPeerFromTCPAddress(ctx context.Context, h *Host, addr string) (*Peer, error) {
	conn, err := h.dialer.DialContext(ctx, "tcp4", addr)
	if err != nil {
		return nil, err
	}

	peer := Peer{
		ID:       rand.Uint64(),
		State:    PeerStateBeforeHandshake,
		protocol: &LevinProtocol{conn},
	}

	return &peer, nil
}

func (p *Peer) Handshake(h *Host) (*HandshakeResponse, error) {
	if p.State != PeerStateBeforeHandshake {
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

	p.State = PeerStateSynchronizing
	p.Version = res.NodeData.Version
	p.Height = res.PayloadData.CurrentHeight

	return &res, nil
}

func (p *Peer) Ping() (*PingResponse, error) {
	req := PingRequest{}
	res := PingResponse{}

	if err := p.protocol.Invoke(CommandPing, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
