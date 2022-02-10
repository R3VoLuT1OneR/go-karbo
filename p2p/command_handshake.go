package p2p

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"go.uber.org/zap"
	"strconv"
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

var (
	ErrHandshakeWrongNetwork    = errors.New("wrong network connection")
	ErrHandshakeNotIncoming     = errors.New("handshake not from incoming connection")
	ErrHandshakeHasID           = errors.New("ID for the connecting peer already set")
	ErrHandshakeProcessSyncData = errors.New("failed to process sync data")
)

// NewHandshakeRequest returns new struct to be sent as handshake request command to new peer.
func NewHandshakeRequest(bc *cryptonote.BlockChain) HandshakeRequest {
	return HandshakeRequest{
		NodeData:    newBasicNodeData(bc.Network),
		PayloadData: *newSyncData(bc),
	}
}

// NewHandshakeResponse returns new struct to make response for handshake request
func NewHandshakeResponse(bc *cryptonote.BlockChain, peerEntries []PeerEntry) HandshakeResponse {
	return HandshakeResponse{
		NodeData:    newBasicNodeData(bc.Network),
		PayloadData: *newSyncData(bc),
		Peers:       peerEntries,
	}
}

func (n *Node) HandleHandshake(p *Peer, req HandshakeRequest) error {
	p.procMutex.Lock()
	defer p.procMutex.Unlock()

	// Update or set remove peer version
	p.SetVersion(req.NodeData.Version)

	// TODO: Verify that remote IP allowed to connect to our node

	if req.NodeData.NetworkID != n.Config.Network.NetworkID {
		err := ErrHandshakeWrongNetwork
		p.logger.Error(err, zap.String("peerID", req.NodeData.NetworkID.String()))
		p.Shutdown()
		return err
	}

	if !p.isIncoming {
		err := ErrHandshakeNotIncoming
		p.logger.Error(err)
		p.Shutdown()
		return err
	}

	if p.ID != 0 {
		err := ErrHandshakeHasID
		p.logger.Error(err, zap.String("peerID", strconv.FormatUint(req.NodeData.PeerID, 10)))
		p.Shutdown()
		return err
	}

	p.SetID(req.NodeData.PeerID)

	if err := p.processSyncData(req.PayloadData, true); err != nil {
		err := ErrHandshakeProcessSyncData
		p.Shutdown()
		return err
	}

	// TODO: Send ping and make sure we can connect to the peer and add it to the white list.

	p.logger.Debug("handshake received")
	return nil
}
