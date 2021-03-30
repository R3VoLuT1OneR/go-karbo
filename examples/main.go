package main

import (
	"context"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/p2p"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main()  {
	mainnet := config.MainNet()

	db, err := cryptonote.NewBadgerDB()
	if err != nil {
		panic(err)
	}

	core, err := cryptonote.NewCore(mainnet, db)
	if err != nil {
		panic(err)
	}

	ctx := interruptListener()
	cfg := p2p.HostConfig{
		BindAddr: "127.0.0.1:32447",
		Network: mainnet,
	}

	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.TraceLevel

	host := p2p.NewNode(core, cfg, logger)

	fmt.Println("Server started.")

	if err := host.Run(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Server stopped.")
}

func interruptListener() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, os.Interrupt)

		select {
		case sig := <-interruptChannel:
			fmt.Printf("Received signal (%s). Shutting down...\n", sig)
		}

		cancel()
	}()

	return ctx
}
