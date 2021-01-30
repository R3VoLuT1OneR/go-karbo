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
	CommandPoolBase 			= 1000
	CommandHandshake 			= CommandPoolBase + 1
	CommandTimedSync 			= CommandPoolBase + 2
	CommandPing 				= CommandPoolBase + 2
	CommandRequestStatInfo 		= CommandPoolBase + 4
	CommandRequestNetworkState 	= CommandPoolBase + 5
	CommandRequestPeerID 		= CommandPoolBase + 6
)

type BasicNodeData struct {
	NetworkID uuid.UUID `binary:"network_id"`
	Version   uint8     `binary:"version"`
	PeerID    uint64    `binary:"peer_id"`
	LocalTime uint64    `binary:"local_time"`
	MyPort    uint32    `binary:"my_port"`
}

type SyncData struct {
	CurrentHeight uint32          `binary:"current_height"`
	TopBlockHash  cryptonote.Hash `binary:"top_id"`
}

type NetworkAddress struct {
	IP uint32
	Port uint32
}

func (na *NetworkAddress) String() string {
	return fmt.Sprintf("%s:%d", ipv4.ToDots(na.IP), na.Port)
}

type PeerEntry struct {
	Address NetworkAddress
	ID uint64
	LastSeen uint64
}

type HandshakeRequest struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
}

type HandshakeResponse struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
	Peers       []PeerEntry   `binary:"local_peerlist"`
}

type TimedSyncRequest struct {
	PayloadData SyncData      `binary:"payload_data"`
}

type TimedSyncResponse struct {
	LocalTime 	uint64 			`binary:"local_time"`
	PayloadData SyncData      	`binary:"payload_data"`
	Peers       []PeerEntry   	`binary:"local_peerlist"`
}

type PingRequest struct {}

type PingResponse struct {
	Status string `binary:"status"`
	PeerID uint64 `binary:"peer_id"`
}

var mapCommandStructs = map[uint32]interface{}{
	CommandTimedSync: TimedSyncRequest{},
}

func NewHandshakeRequest(network *config.Network) (*HandshakeRequest, error) {
	syncData, err := prepareSyncData(network)
	if err != nil {
		return nil, err
	}

	r := HandshakeRequest{
		NodeData: BasicNodeData{
			NetworkID: network.NetworkID,
			Version:   network.P2PCurrentVersion,
			LocalTime: uint64(time.Now().Unix()),
			PeerID:    uint64(rand.Int63()),
			MyPort:    32347,
		},
		PayloadData: *syncData,
	}

	return &r, nil
}

func newTimedSyncResponse(n *config.Network) (*TimedSyncResponse, error) {
	syncData, err := prepareSyncData(n)
	if err != nil {
		return nil, err
	}

	var peers []PeerEntry
	return &TimedSyncResponse{uint64(time.Now().Unix()), *syncData, peers}, nil
}

func prepareSyncData(network *config.Network) (*SyncData, error) {
	// TODO: Top block must be fetched from blockchain storage
	topBlock, err := cryptonote.GenerateGenesisBlock(network)
	if err != nil {
		return nil, err
	}

	hash, err := topBlock.Hash()
	if err != nil {
		return nil, err
	}

	return &SyncData{uint32(0), *hash}, nil
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
