package p2p

const PingCommandID = CommandPoolBase + 3

type PingRequest struct {}

type PingResponse struct {
	status string
	peer_id uint64
}
