package main

import (
	"context"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/p2p"
	"net"
)

var mainnet = config.MainNet()

func main()  {

	ctx := context.Background()

	ba, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:32347")
	if err != nil {
		panic(err)
	}

	cfg := p2p.HostConfig{
		BindAddr: ba,
		Network: mainnet,
	}

	host := p2p.NewHost(cfg)

	err = host.Start(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("host", host)
}


