package p2p

import "github.com/r3volut1oner/go-karbo/cryptonote"

type HandshakeRequest struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
}

type HandshakeResponse struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
	Peers       []PeerEntry   `binary:"local_peerlist,binary"`
}

func NewHandshakeRequest(bc *cryptonote.BlockChain) (*HandshakeRequest, error) {
	syncData, err := newSyncData(bc)
	if err != nil {
		return nil, err
	}

	nodeData, err := newBasicNodeData(bc.Network)
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

	nodeData, err := newBasicNodeData(h.Blockchain.Network)
	if err != nil {
		return nil, err
	}

	payloadData, err := newSyncData(h.Blockchain)
	if err != nil {
		return nil, err
	}

	return &HandshakeResponse{
		NodeData:    nodeData,
		PayloadData: *payloadData,
		Peers:       peerList,
	}, nil
}
