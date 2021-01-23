package main

import (
	"bufio"
	"context"
	"fmt"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	realpeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/r3volut1oner/go-karbo/config"
)

var mainnet = config.MainNet()

func SeedHostToMultiAddr(addr string) (ma.Multiaddr, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	multAddr, err := manet.FromNetAddr(netAddr)
	if err != nil {
		return nil, err
	}

	_, SeedKey, err := p2pcrypto.GenerateKeyPair(p2pcrypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}

	SeedId, err := peerstore.IDFromPublicKey(SeedKey)
	if err != nil {
		return nil, err
	}

	multAddrSeedId, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", SeedId));
	if err != nil {
		return nil, err
	}

	return multAddr.Encapsulate(multAddrSeedId), nil
}


func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func main() {

	// create a background context (i.e. one that never cancels)
	ctx := context.Background()

	// start a libp2p srv that listens on TCP port 2000 on the IPv4
	// loopback interface
	srv, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)

	if err != nil {
		panic(err)
	}

	// Peer ID
	peerInfo := peerstore.AddrInfo{
		ID:    srv.ID(),
		Addrs: srv.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p srv address:", addrs[0])

	// print the srv's listening addresses
	fmt.Println("Listen addresses:", srv.Addrs())

	var seeds []ma.Multiaddr

	for i := 0; i < len(mainnet.SeedNodes); i++ {
		seedHost := mainnet.SeedNodes[i]
		seedMultiAddr, err := SeedHostToMultiAddr(seedHost)

		if err != nil {
			panic(err)
		}

		seeds = append(seeds, seedMultiAddr)
	}

	fmt.Println("Seeds", seeds)

	srv.SetStreamHandler("/karbo/0.0.1", func(stream network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		go readData(rw)
	})

	for i := 0; i < len(seeds); i++ {
		addr := seeds[i]

		fmt.Println("AddrInfoFromP2pAddr", addr)

		peer, err := peerstore.AddrInfoFromP2pAddr(seeds[i])
		if err != nil {
			panic(err)
		}

		srv.Peerstore().AddAddr(peer.ID, peer.Addrs[0], realpeerstore.PermanentAddrTTL)

		fmt.Println("Peer", peer)

		if err := srv.Connect(ctx, *peer); err != nil {
			panic(err)
		}
	}

	// Configure ping protocol
	//pingService := &ping.PingService{Host: srv}
	//srv.SetStreamHandler(ping.ID, pingService.PingHandler)
	//
	//if len(os.Args) > 1 {
	//	addr, err := multiaddr.NewMultiaddr(os.Args[1])
	//	if err != nil {
	//		panic(err)
	//	}
	//	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if err := srv.Connect(ctx, *peer); err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("Sending 5 ping request to", addr)
	//	ch := pingService.Ping(ctx, peer.ID)
	//	for i := 0; i < 5; i++ {
	//		res := <-ch
	//		fmt.Println("go ping response!", "RTT:", res.RTT)
	//	}
	//} else {
	//	// Listening for termination signals
	//	termChannel := make(chan os.Signal, 1)
	//	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM)
	//	<-termChannel
	//	fmt.Println("Received signal, shutting down...")
	//}

	// shut the srv down
	if err := srv.Close(); err != nil {
		panic(err)
	}
}
