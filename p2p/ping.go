package p2p

const PingCommandID = CommandPoolBase + 3

type PingRequest struct {}

type PingResponse struct {
	Status string `binary:"status"`
	PeerID uint64 `binary:"peer_id"`
}
