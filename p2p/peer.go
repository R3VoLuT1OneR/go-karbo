package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/cryptonote"
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
	// TODO: Make sure we generate local peer ID and updating external IDs
	ID uint64

	core    *cryptonote.Core
	version byte
	state   byte

	address *net.TCPAddr

	protocol *LevinProtocol

	remoteHeight uint32
	lastResponseHeight uint32

	neededBlocks 	cryptonote.HashList
	requestedBlocks cryptonote.HashList
}

func NewPeerFromTCPAddress(ctx context.Context, h *Node, addr string) (*Peer, error) {
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
		core: h.Core,
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

		for {
			nb := p.neededBlocks[0]

			if !(checkHavingBlocks && p.core.HasBlock(&nb)) {
				n.Blocks = append(n.Blocks, nb)
				p.requestedBlocks = append(p.requestedBlocks, nb)
			}

			p.neededBlocks = p.neededBlocks[1:]

			if len(n.Blocks) >= MaxBlockSynchronization || len(p.neededBlocks) == 0 {
				break
			}
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