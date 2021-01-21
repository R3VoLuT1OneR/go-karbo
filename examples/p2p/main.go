package main

import (
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"net"
)

func main()  {

	params := config.MainNetParams()
	seeds := params.SeedNodes

	for _, addr := range(seeds) {
		fmt.Println("Seed address:", addr)

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		var r []byte

		l, err := conn.Read(r)
		if err != nil {
			panic(err)
		}

		fmt.Println("Read", l, r)
	}

}
