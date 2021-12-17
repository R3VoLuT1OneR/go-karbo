package p2p

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/signalsciences/ipv4"
	"math/rand"
	"reflect"
	"time"
)

const (
	CommandPoolBase            = 1000
	CommandHandshake           = CommandPoolBase + 1
	CommandTimedSync           = CommandPoolBase + 2
	CommandPing                = CommandPoolBase + 3
	CommandRequestStatInfo     = CommandPoolBase + 4
	CommandRequestNetworkState = CommandPoolBase + 5
	CommandRequestPeerID       = CommandPoolBase + 6
)

type BasicNodeData struct {
	NetworkID uuid.UUID `binary:"network_id,binary"`
	Version   uint8     `binary:"version"`
	PeerID    uint64    `binary:"peer_id"`
	LocalTime uint64    `binary:"local_time"`
	MyPort    uint32    `binary:"my_port"`
}

type SyncData struct {
	CurrentHeight uint32          `binary:"current_height"`
	TopBlockHash  cryptonote.Hash `binary:"top_id,binary"`
}

type NetworkAddress struct {
	IP   uint32
	Port uint32
}

func (na *NetworkAddress) String() string {
	return fmt.Sprintf("%s:%d", ipv4.ToDots(na.IP), na.Port)
}

type PeerEntry struct {
	Address  NetworkAddress
	ID       uint64
	LastSeen uint64
}

func (pe *PeerEntry) FromPeer(p *Peer) error {
	IP, err := ipv4.FromNetIP(p.address.IP)
	if err != nil {
		return err
	}

	pe.ID = p.ID
	pe.Address = NetworkAddress{
		IP:   IP,
		Port: uint32(p.address.Port),
	}

	// TODO: Get real last seen
	pe.LastSeen = uint64(time.Now().Unix())

	return nil
}

type HandshakeRequest struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
}

type HandshakeResponse struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
	Peers       []PeerEntry   `binary:"local_peerlist,binary"`
}

type TimedSyncRequest struct {
	PayloadData SyncData `binary:"payload_data"`
}

type TimedSyncResponse struct {
	LocalTime   uint64      `binary:"local_time"`
	PayloadData SyncData    `binary:"payload_data"`
	Peers       []PeerEntry `binary:"local_peerlist,binary"`
}

type PingRequest struct{}

type PingResponse struct {
	Status string `binary:"status"`
	PeerID uint64 `binary:"peer_id"`
}

var mapCommandStructs = map[uint32]interface{}{
	CommandHandshake: HandshakeRequest{},
	CommandTimedSync: TimedSyncRequest{},
}

func NewHandshakeRequest(core *cryptonote.Core) (*HandshakeRequest, error) {
	syncData, err := newSyncData(core)
	if err != nil {
		return nil, err
	}

	nodeData, err := newBasicNodeData(core.Network)
	if err != nil {
		return nil, err
	}

	r := HandshakeRequest{
		NodeData:    nodeData,
		PayloadData: *syncData,
	}

	return &r, nil
}

func NewHandshakeResponse(h *Node) (*HandshakeResponse, error) {
	peerList, err := newPeerEntryList(h)
	if err != nil {
		return nil, err
	}

	nodeData, err := newBasicNodeData(h.Core.Network)
	if err != nil {
		return nil, err
	}

	payloadData, err := newSyncData(h.Core)
	if err != nil {
		return nil, err
	}

	return &HandshakeResponse{
		NodeData:    nodeData,
		PayloadData: *payloadData,
		Peers:       peerList,
	}, nil
}

func newTimedSyncResponse(h *Node) (*TimedSyncResponse, error) {
	syncData, err := newSyncData(h.Core)
	if err != nil {
		return nil, err
	}

	peerList, err := newPeerEntryList(h)
	if err != nil {
		return nil, err
	}

	return &TimedSyncResponse{
		LocalTime:   uint64(time.Now().Unix()),
		PayloadData: *syncData,
		Peers:       peerList,
	}, nil
}

func newBasicNodeData(n *config.Network) (BasicNodeData, error) {
	return BasicNodeData{
		NetworkID: n.NetworkID,
		Version:   n.P2PCurrentVersion,
		LocalTime: uint64(time.Now().Unix()),
		PeerID:    uint64(rand.Int63()),
		MyPort:    32347,
	}, nil
}

func newPeerEntryList(h *Node) ([]PeerEntry, error) {
	var peers []PeerEntry

	for _, p := range h.ps.white.peers {
		var pe PeerEntry
		if err := pe.FromPeer(p); err != nil {
			return nil, err
		}

		peers = append(peers, pe)
	}

	for _, p := range h.ps.grey.peers {
		var pe PeerEntry
		if err := pe.FromPeer(p); err != nil {
			return nil, err
		}

		peers = append(peers, pe)
	}

	return peers, nil
}

func newSyncData(core *cryptonote.Core) (*SyncData, error) {
	topBlock, height, err := core.TopBlock()
	if err != nil {
		return nil, err
	}

	hash := topBlock.Hash()

	return &SyncData{height, *hash}, nil
}

func parseCommand(lc *LevinCommand) (interface{}, error) {
	if s, ok := mapCommandStructs[lc.Command]; ok {
		command := reflect.New(reflect.TypeOf(s))
		if err := binary.Unmarshal(lc.Payload, command.Interface()); err != nil {
			return nil, err
		}

		return command.Elem().Interface(), nil
	}

	return nil, errors.New(fmt.Sprintf("unknown command: %d", lc.Command))
}
