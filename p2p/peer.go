package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/signalsciences/ipv4"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
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

type NetworkAddress struct {
	IP   uint32
	Port uint32
}

type Peer struct {
	// TODO: Make sure we generate local peer ID and updating external IDs
	ID uint64

	// logger to be used for logging any peer events
	logger *zap.SugaredLogger

	// isIncoming flags if peer is got from incoming connection
	isIncoming bool

	version byte
	state   byte

	address NetworkAddress

	protocol *LevinProtocol

	remoteHeight       uint32
	lastResponseHeight uint32

	neededBlocks    crypto.HashList
	requestedBlocks crypto.HashList

	// procMutex blocked when peer processing some command or notification
	procMutex sync.Mutex

	// struct mutex is used for blocking peer attribute updates
	sync.RWMutex
}

func NetworkAddressFromTCPAddr(addr *net.TCPAddr) NetworkAddress {
	IP, _ := ipv4.FromNetIP(addr.IP)

	return NetworkAddress{
		IP:   IP,
		Port: uint32(addr.Port),
	}
}

// NewPeerFromTCPAddress returns new peer from the IP address. Used for creating seed peers.
func NewPeerFromTCPAddress(ctx context.Context, n *Node, addr string) (*Peer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := n.dialer.DialContext(ctx, "tcp4", addr)
	if err != nil {
		return nil, err
	}

	address := NetworkAddressFromTCPAddr(tcpAddr)

	return NewPeer(n.logger, &LevinProtocol{conn}, address, false), nil
}

// NewPeerFromIncomingConnection returns new seed from some incoming connection.
func NewPeerFromIncomingConnection(n *Node, conn *net.TCPConn) *Peer {
	address := NetworkAddressFromTCPAddr(conn.RemoteAddr().(*net.TCPAddr))

	return NewPeer(n.logger, &LevinProtocol{conn}, address, true)
}

func NewPeer(logger *zap.SugaredLogger, protocol *LevinProtocol, address NetworkAddress, isIncoming bool) *Peer {
	return &Peer{
		protocol:   protocol,
		address:    address,
		isIncoming: isIncoming,
		logger: logger.With(
			zap.String("address", address.String()),
			zap.Bool("isIncoming", isIncoming),
		),
	}
}

func (p *Peer) SetID(ID uint64) {
	p.Lock()
	p.logger = p.logger.With(zap.String("ID", strconv.FormatUint(ID, 10)))
	p.ID = ID
	p.Unlock()
}

func (p *Peer) SetVersion(version byte) {
	p.Lock()
	p.logger = p.logger.With(zap.String("version", string(version)))
	p.version = version
	p.Unlock()
}

func (p *Peer) PeerEntry() PeerEntry {
	return PeerEntry{
		ID:       p.ID,
		Address:  p.address,
		LastSeen: uint64(time.Now().Unix()),
	}
}

func (p *Peer) Shutdown() {
	p.logger.Debug("shutdown request received")
	p.state = PeerStateShutdown
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s", p.address.String())
}

func (p *Peer) handshake(n *Node) (*HandshakeResponse, error) {
	if p.state != PeerStateBeforeHandshake {
		return nil, errors.New("state is not before handshake")
	}

	var res HandshakeResponse
	if err := p.protocol.Invoke(CommandHandshake, NewHandshakeRequest(n.Blockchain), &res); err != nil {
		return nil, err
	}

	if n.Config.Network.NetworkID != res.NodeData.NetworkID {
		return nil, errors.New("wrong network id received")
	}

	if n.Config.Network.P2PMinimumVersion < res.NodeData.Version {
		return nil, errors.New("node data version not match minimal")
	}

	if err := n.processSyncData(p, res.PayloadData, true); err != nil {
		return nil, err
	}

	p.version = res.NodeData.Version
	p.ID = res.NodeData.PeerID

	// TODO: Handle new peerlist

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

func (p *Peer) requestChain(bc *cryptonote.BlockChain) error {
	n, err := NewRequestChain(bc)
	if err != nil {
		return err
	}

	p.logger.Debugf("[%s] request chain %d (%d blocks) ", p, p.lastResponseHeight, len(n.Blocks))

	if err := p.protocol.Notify(NotificationRequestChainID, *n); err != nil {
		return err
	}

	return nil
}

func (p *Peer) requestMissingBlocks(n *Node, checkHavingBlocks bool) error {
	if len(p.neededBlocks) > 0 {
		neededBlocks := p.neededBlocks
		requestBlocks := crypto.HashList{}

		for len(neededBlocks) > 0 && len(requestBlocks) < MaxBlockSynchronization {
			nb := neededBlocks[0]

			haveBlock := n.Blockchain.HaveBlock(&nb)

			if !(checkHavingBlocks && haveBlock) {
				requestBlocks = append(requestBlocks, nb)
			}

			neededBlocks = neededBlocks[1:]
		}

		if len(requestBlocks) > 0 {
			n := NotificationRequestGetObjects{}
			n.Blocks = requestBlocks

			if err := p.protocol.Notify(NotificationRequestGetObjectsID, n); err != nil {
				return err
			}

			p.requestedBlocks = append(p.requestedBlocks, requestBlocks...)
		}

		p.neededBlocks = neededBlocks
	} else if p.lastResponseHeight < (p.remoteHeight - 1) {
		if err := p.requestChain(n.Blockchain); err != nil {
			return err
		}
	} else {
		if p.lastResponseHeight == p.remoteHeight-1 && len(p.neededBlocks) != 0 && len(p.requestedBlocks) != 0 {
			return errors.New(fmt.Sprintf(
				"request missing blocks final condition failed: \n"+
					"response height: %d\n"+
					"remote blockchain height: %d\n"+
					"needed objects size: %d\n"+
					"requested objects size: %d",
				p.lastResponseHeight,
				p.remoteHeight,
				len(p.neededBlocks),
				len(p.requestedBlocks),
			))
		}

		// TODO: Request missing pool transactions
		// src/CryptoNoteProtocol/CryptoNoteProtocolHandler.cpp:907

		p.state = PeerStateNormal
		p.logger.Debugf("[%s] syncronized", p)

		// TODO: On connection synchronized
		// src/CryptoNoteProtocol/CryptoNoteProtocolHandler.cpp:911
	}

	return nil
}

func (na *NetworkAddress) String() string {
	return fmt.Sprintf("%s:%d", ipv4.ToDots(na.IP), na.Port)
}
