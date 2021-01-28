package p2p

import (
	"github.com/google/uuid"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"math/rand"
	"time"
)

type BasicNodeData struct {
	NetworkID 	uuid.UUID 	`binary:"network_id"`
	Version 	uint8 		`binary:"version"`
	PeerId 		uint64 		`binary:"peer_id"`
	LocalTime 	uint64 		`binary:"local_time"`
	MyPort 		uint32 		`binary:"my_port"`
}

type SyncData struct {
	CurrentHeight uint32          `binary:"current_height"`
	TopBlockHash  cryptonote.Hash `binary:"top_id"`
}

type NetworkAddress struct {
	IP uint32
	Port uint32
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

func NewHandshakeRequest(network *config.Network) (*HandshakeRequest, error) {
	// TODO: Top block must be fetched from blockchain storage
	topBlock, err := cryptonote.GenerateGenesisBlock(network)
	if err != nil {
		return nil, err
	}

	hash, err := topBlock.Hash()
	if err != nil {
		return nil, err
	}

	// TODO: Must be fetched from blockchain storage
	height := 0

	r := HandshakeRequest{
		NodeData: BasicNodeData{
			NetworkID: network.NetworkID,
			Version: network.P2PCurrentVersion,
			LocalTime: uint64(time.Now().Unix()),
			PeerId: uint64(rand.Int63()),
			MyPort: 32347,
		},
		PayloadData: SyncData{
			CurrentHeight: uint32(height),
			TopBlockHash: *hash,
		},
	}

	return &r, nil
}
