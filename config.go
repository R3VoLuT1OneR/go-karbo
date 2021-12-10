package main

import (
	flags "github.com/jessevdk/go-flags"
)

type config struct {
	ShowVersion 		bool `short:"V" long:"version" description:"Display version information and exit"`
}

func newConfigParser(cfg *config, so *serviceOptions)
