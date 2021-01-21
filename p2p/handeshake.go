package p2p

import (
	"github.com/google/uuid"
	"github.com/r3volut1oner/go-karbo/cryptonote/block"
)

const HandshakeCommandID = CommandPoolBase + 1

type BasicNodeData struct {
	NetworkID uuid.UUID `binary:"network_id"`
	Version uint8 `binary:"version"`
	PeerId uint64 `binary:"peer_id"`
	LocalTime uint64 `binary:"local_time"`
	MyPort uint32 `binary:"my_port"`
}

type PayloadData struct {
	CurrentHeight uint32 `binary:"current_height"`
	TopBlockHash block.HashBytes `binary:"top_id"`
}

type NetworkAddress struct {
	IP uint32
	Port uint32
}

type Peer struct {
	Address NetworkAddress
	ID uint64
	LastSeen uint64
}

type HandshakeRequest struct {
	NodeData BasicNodeData `binary:"node_data"`
	PayloadData PayloadData `binary:"payload_data"`
}

type HandshakeResponse struct {
	NodeData BasicNodeData `binary:"node_data"`
	PayloadData PayloadData `binary:"payload_data"`
	Peers []Peer `binary:"local_peerlist"`
}

