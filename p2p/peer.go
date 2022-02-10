package p2p

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/signalsciences/ipv4"
	"go.uber.org/zap"
	"io/ioutil"
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

func (p *Peer) handleResponseGetObjects(bc *cryptonote.BlockChain, nt NotificationResponseGetObjects) error {

	p.logger.Debugf("[%s] response to get objects", p)

	if len(nt.Blocks) == 0 {
		p.Shutdown()
		return errors.New(fmt.Sprintf("[%s] got zero blocks on get objects", p))
	}

	if p.lastResponseHeight > nt.CurrentBlockchainHeight {
		p.Shutdown()
		return errors.New(fmt.Sprintf(
			"[%s] got wrong currentBlockchainHeight = %d, current = %d", p,
			nt.CurrentBlockchainHeight,
			p.lastResponseHeight,
		))
	}

	// TODO: Implement P2P Node observable height (max observed height) and update it if new observed height is found.

	p.remoteHeight = nt.CurrentBlockchainHeight

	newObjects := map[*cryptonote.Block][][]byte{}
	for i, rawBlock := range nt.Blocks {
		block := cryptonote.Block{}
		rawBlockReader := bytes.NewReader(rawBlock.Block)
		if err := block.Deserialize(rawBlockReader); err != nil {
			p.Shutdown()
			// TODO: Remove this. It is debug only.
			height := bc.Height()
			blockHeight := height + uint32(i)
			_ = ioutil.WriteFile(fmt.Sprintf("./block_%d.dat", blockHeight), rawBlock.Block, 0644)
			return errors.New(
				fmt.Sprintf("[%s] (%d) failed to convert raw block (%d): %s", p, i, blockHeight, err),
			)
		}

		hash := block.Hash()
		if !p.requestedBlocks.Has(hash) {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] got not requested block #%d '%s'", p, i, hash.String()))
		}

		if len(block.TransactionsHashes) != len(rawBlock.Transactions) {
			p.Shutdown()
			return errors.New(fmt.Sprintf(
				"[%s] got wrong block transactions size. block: %s block tx: %d raw tx: %d",
				p, hash.String(), len(block.TransactionsHashes), len(rawBlock.Transactions),
			))
		}

		p.requestedBlocks.Remove(hash)
		newObjects[&block] = rawBlock.Transactions
	}

	if len(p.requestedBlocks) > 0 {
		p.Shutdown()
		return errors.New(fmt.Sprintf(
			"[%s] got not all requested objectes, missing %d", p, len(p.requestedBlocks),
		))
	}

	if err := p.processNewObjects(bc, newObjects); err != nil {
		return err
	}

	height := bc.Height()
	p.logger.Infof("process block, total height: %d", height)

	return p.requestMissingBlocks(bc, true)
}

func (p *Peer) processNewObjects(bc *cryptonote.BlockChain, objects map[*cryptonote.Block][][]byte) error {
	for block, transactions := range objects {
		if err := bc.AddBlock(block, transactions); err != nil {
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

func (p *Peer) handshake(h *Node) (*HandshakeResponse, error) {
	if p.state != PeerStateBeforeHandshake {
		return nil, errors.New("state is not before handshake")
	}

	var res HandshakeResponse
	if err := p.protocol.Invoke(CommandHandshake, NewHandshakeRequest(h.Blockchain), &res); err != nil {
		return nil, err
	}

	if h.Config.Network.NetworkID != res.NodeData.NetworkID {
		return nil, errors.New("wrong network id received")
	}

	if h.Config.Network.P2PMinimumVersion < res.NodeData.Version {
		return nil, errors.New("node data version not match minimal")
	}

	if err := p.processSyncData(res.PayloadData, true); err != nil {
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
	n, err := newRequestChain(bc)
	if err != nil {
		return err
	}

	p.logger.Debugf("[%s] request chain %d (%d blocks) ", p, p.lastResponseHeight, len(n.Blocks))

	if err := p.protocol.Notify(NotificationRequestChainID, *n); err != nil {
		return err
	}

	return nil
}

func (p *Peer) processSyncData(data SyncData, initial bool) error {
	// TODO: Implement sync data
	if p.state == PeerStateBeforeHandshake && !initial {
		return nil
	}

	if p.state == PeerStateSynchronizing {
	} else {
		p.state = PeerStateSyncRequired
	}

	p.remoteHeight = data.CurrentHeight

	return nil
}

func (p *Peer) requestMissingBlocks(bc *cryptonote.BlockChain, checkHavingBlocks bool) error {
	if len(p.neededBlocks) > 0 {
		neededBlocks := p.neededBlocks
		requestBlocks := crypto.HashList{}

		for len(neededBlocks) > 0 && len(requestBlocks) < MaxBlockSynchronization {
			nb := neededBlocks[0]

			haveBlock := bc.HaveBlock(&nb)

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
		if err := p.requestChain(bc); err != nil {
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
