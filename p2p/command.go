package p2p

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"math/rand"
	"reflect"
	"time"
)

const (
	CommandPoolBase            = 1000                // 1000
	CommandHandshake           = CommandPoolBase + 1 // 1001
	CommandTimedSync           = CommandPoolBase + 2 // 1002
	CommandPing                = CommandPoolBase + 3 // 1003
	CommandRequestStatInfo     = CommandPoolBase + 4 // 1004
	CommandRequestNetworkState = CommandPoolBase + 5 // 1005
	CommandRequestPeerID       = CommandPoolBase + 6 // 1006
)

type BasicNodeData struct {
	NetworkID uuid.UUID `binary:"network_id,binary"`
	Version   uint8     `binary:"version"`
	PeerID    uint64    `binary:"peer_id"`
	LocalTime uint64    `binary:"local_time"`
	MyPort    uint32    `binary:"my_port"`
}

type SyncData struct {
	CurrentHeight uint32      `binary:"current_height"`
	TopBlockHash  crypto.Hash `binary:"top_id,binary"`
}

type PeerEntry struct {
	Address  NetworkAddress
	ID       uint64
	LastSeen uint64
}

func (pe *PeerEntry) FromPeer(p *Peer) error {
	pe.ID = p.ID
	pe.Address = p.address

	// TODO: Get real last seen
	pe.LastSeen = uint64(time.Now().Unix())

	return nil
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

func newTimedSyncResponse(h *Node) (*TimedSyncResponse, error) {
	syncData, err := newSyncData(h.Blockchain)
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

func newSyncData(bc *cryptonote.BlockChain) (*SyncData, error) {
	topBlock := bc.TopBlock()

	hash := topBlock.Hash()
	height := topBlock.Index() + 1

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
