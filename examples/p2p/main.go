package main

import (
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/r3volut1oner/go-karbo/p2p"
	"net"
)

func main()  {

	network := config.MainNet()
	seeds := network.SeedNodes

	for _, addr := range(seeds) {
		fmt.Println("Seed address:", addr)

		req, err := p2p.NewHandshakeRequest(network)
		if err != nil {
			panic(err)
		}

		reqBytes, err := binary.Marshal(*req)
		if err != nil {
			panic(err)
		}

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		tcpc := conn.(*net.TCPConn)
		if err := tcpc.SetKeepAlive(true); err != nil {
			panic(err)
		}

		levin := p2p.NewLevinProtocol(conn)
		if _, err := levin.WriteCommand(p2p.CommandHandshake, reqBytes, true); err != nil {
			panic(err)
		}

		command, err := levin.ReadCommand()
		if err != nil {
			panic(err)
		}

		fmt.Println("command", command.Command)

		if command.Command == p2p.CommandHandshake {
			var rsp p2p.HandshakeResponse

			if err := binary.Unmarshal(command.Payload, &rsp); err != nil {
				panic(err)
			}

			fmt.Println("rsp", rsp)
		}
	}

}
