package main

import (
	"fmt"
	"os"

	"github.com/r3volut1oner/go-karbo/config"
)

var mainNetParams = config.MainNet()

func main()  {
	_, err := fmt.Fprintf(os.Stdin, "Seed Nodes: %v\n", mainNetParams.SeedNodes)

	if err != nil {
		panic(err)
	}
}


