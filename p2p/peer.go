package p2p

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
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
	Height  uint32
	state   byte

	address *net.TCPAddr

	protocol *LevinProtocol
}

func NewPeerFromTCPAddress(ctx context.Context, h *Host, addr string) (*Peer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := h.dialer.DialContext(ctx, "tcp4", addr)
	if err != nil {
		return nil, err
	}

	peer := Peer{
		ID:       rand.Uint64(),
		protocol: &LevinProtocol{conn},
		address: tcpAddr,
	}

	return &peer, nil
}

func NewPeerFromIncomingConnection(conn net.Conn) *Peer {
	return &Peer{
		protocol: &LevinProtocol{conn},
	}
}

func (p *Peer) handshake(h *Host) (*HandshakeResponse, error) {
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
	p.Version = res.NodeData.Version
	p.Height = res.PayloadData.CurrentHeight

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

func (p *Peer) requestChain(h *Host) error {
	requestChain, err := newRequestChain(h.Config.Network)
	if err != nil {
		return err
	}

	fmt.Println("requestChain", requestChain)

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
