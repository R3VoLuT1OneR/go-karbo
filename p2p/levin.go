package p2p

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	// const uint64_t LEVIN_SIGNATURE = 0x0101010101012101LL;  //Bender's nightmare
	LevinSignature uint64 = 0x0101010101012101  //Bender's nightmare

	// const uint32_t LEVIN_PACKET_REQUEST = 0x00000001;
	LevinPacketRequest uint32 = 0x00000001

	// const uint32_t LEVIN_PACKET_RESPONSE = 0x00000002;
	LevinPacketResponse uint32 = 0x00000002

	// const uint32_t LEVIN_DEFAULT_MAX_PACKET_SIZE = 100000000;      //100MB by default
	LevinMaxPacketSize uint32 = 100000000

	// const uint32_t LEVIN_PROTOCOL_VER_1 = 1;
	LevinProtocolVersion1 uint32 = 1

	LevinHeadSize = 33
)

type LevinProtocol struct {
	Conn net.Conn
}

type LevinCommand struct {
	Command uint32
	IsNotify bool
	IsResponse bool
	Payload []byte
}

type bucketHead struct {
	Signature uint64
	BodySize uint64
	HaveToReturnData bool
	Command uint32
	ReturnCode int32
	Flags uint32
	ProtocolVersion uint32
}

func NewLevinProtocol(conn net.Conn) LevinProtocol {
	return LevinProtocol{conn}
}

func (protocol *LevinProtocol) WriteCommand(command uint32, payload []byte, HaveToReturnData bool) (n int, err error) {
	return protocol.send(command, payload, HaveToReturnData, LevinPacketRequest, 0)
}

func (protocol *LevinProtocol) WriteReply(command uint32, payload []byte, returnCode int32) (n int, err error){
	return protocol.send(command, payload, false, LevinPacketResponse, returnCode)
}

func (protocol *LevinProtocol) ReadCommand() (*LevinCommand, error) {
	var headBytes [LevinHeadSize]byte
	var head bucketHead

	if _, err := io.ReadFull(protocol.Conn, headBytes[:]); err != nil {
		return nil, err
	}

	head.decode(headBytes)

	if head.Signature != LevinSignature {
		return nil, errors.New("levin signature mismatch")
	}

	if head.BodySize > uint64(LevinMaxPacketSize) {
		return nil, errors.New("levin packet size is too big")
	}

	payload := make([]byte, head.BodySize)

	if _, err := io.ReadFull(protocol.Conn, payload); err != nil {
		return nil, err
	}

	return &LevinCommand{
		Command: head.Command,
		Payload: payload,
		IsNotify: !head.HaveToReturnData,
		IsResponse: (head.Flags & LevinPacketResponse) == LevinPacketResponse,
	}, nil
}

func (protocol *LevinProtocol) send(
	command uint32,
	payload []byte,
	haveToReturnData bool,
	flags uint32,
	returnCode int32,
) (n int, err error) {
	pLen := len(payload)
	head := bucketHead{
		Signature:        LevinSignature,
		BodySize:         uint64(pLen),
		HaveToReturnData: haveToReturnData,
		Command:          command,
		ReturnCode:       returnCode,
		Flags:            flags,
		ProtocolVersion:  LevinProtocolVersion1,
	}

	message := make([]byte, LevinHeadSize+ pLen)

	headBytes := head.encode()
	copy(message[0:33], headBytes[:])
	copy(message[33:], payload)

	return protocol.Conn.Write(message)
}

func (head *bucketHead) encode() [LevinHeadSize]byte {
	var result [LevinHeadSize]byte

	binary.LittleEndian.PutUint64(result[0:8], head.Signature)
	binary.LittleEndian.PutUint64(result[8:16], head.BodySize)

	if head.HaveToReturnData {
		result[16] = 1
	}

	binary.LittleEndian.PutUint32(result[17:21], head.Command)
	binary.LittleEndian.PutUint32(result[21:25], uint32(head.ReturnCode))
	binary.LittleEndian.PutUint32(result[25:29], head.Flags)
	binary.LittleEndian.PutUint32(result[29:33], head.ProtocolVersion)

	return result
}

func (head *bucketHead) decode(b [33]byte) {
	head.Signature = binary.LittleEndian.Uint64(b[0:8])
	head.BodySize = binary.LittleEndian.Uint64(b[8:16])
	head.HaveToReturnData = b[16] != 0
	head.Command = binary.LittleEndian.Uint32(b[17:21])
	head.ReturnCode = int32(binary.LittleEndian.Uint32(b[21:25]))
	head.Flags = binary.LittleEndian.Uint32(b[25:29])
	head.ProtocolVersion = binary.LittleEndian.Uint32(b[29:33])
}
