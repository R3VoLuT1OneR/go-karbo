package p2p

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/cryptonote"
)

type HandshakeRequest struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
}

type HandshakeResponse struct {
	NodeData    BasicNodeData `binary:"node_data"`
	PayloadData SyncData      `binary:"payload_data"`
	Peers       []PeerEntry   `binary:"local_peerlist,binary"`
}

// NewHandshakeRequest returns new struct to be sent as handshake request command to new peer.
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

// NewHandshakeResponse returns new struct to make response for handshake request
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

func HandleHandshake(n *Node, p *Peer, req HandshakeRequest, cmd *LevinCommand) error {
	p.Lock()

	p.version = req.NodeData.Version

	// TODO: Verify that remote IP allowed to connect to our node

	// TODO: Check peer network and rest of the data
	if req.NodeData.NetworkID != n.Config.Network.NetworkID {
		return errors.New("wrong network on handshake")
	}

	// TODO: Send ping and make sure we can connect to the peer and add it to the white list.
	//if err := p.processSyncData(c.(HandshakeRequest).PayloadData, true); err != nil {
	//	return err
	//}

	n.logger.Debugf("[%v] handshake received", p.ID)

	rsp, err := NewHandshakeResponse(n)
	if err != nil {
		return err
	}

	if err := p.protocol.Reply(cmd.Command, *rsp, 1); err != nil {
		return err
	}

	return nil
}
