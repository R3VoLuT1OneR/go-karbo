package p2p

import (
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
)

type Peer struct {
	PeerID uint64
	Protocol *LevinProtocol
}

func (p *Peer) Handshake(h *Host) (*HandshakeResponse, error) {
	req, err := NewHandshakeRequest(h.Config.Network)
	if err != nil {
		return nil, err
	}

	reqBytes, err := binary.Marshal(*req)
	if err != nil {
		return nil, err
	}

	if _, err := p.Protocol.WriteCommand(CommandHandshake, reqBytes, true); err != nil {
		return nil, err
	}

	command, err := p.Protocol.ReadCommand()
	if err != nil {
		return nil, err
	}

	if command.Command != CommandHandshake {
		return nil, errors.New(fmt.Sprintf("wrong command response code: %v", command.Command))
	}

	if !command.IsResponse {
		return nil, errors.New("not response returned")
	}

	var rsp HandshakeResponse
	if err := binary.Unmarshal(command.Payload, &rsp); err != nil {
		return nil, err
	}

	if h.Config.Network.NetworkID != rsp.NodeData.NetworkID {
		return nil, errors.New("wrong network id received")
	}

	if h.Config.Network.P2PMinimumVersion < rsp.NodeData.Version {
		return nil, errors.New("node data version not match minimal")
	}

	p.PeerID = rsp.NodeData.PeerId

	return &rsp, nil
}

func (p *Peer) Ping(h *Host) (*PingResponse, error) {
	req := PingRequest{}

	reqBytes, err := binary.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := p.Protocol.WriteCommand(CommandPing, reqBytes, true); err != nil {
		return nil, err
	}

	command, err := p.Protocol.ReadCommand()
	if err != nil {
		return nil, err
	}

	//if command.Command != CommandPing {
	//	return nil, errors.New(fmt.Sprintf("wrong ping command response code: %v", command.Command))
	//}

	var rsp PingResponse
	if err := binary.Unmarshal(command.Payload, &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}
