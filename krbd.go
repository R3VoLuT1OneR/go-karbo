package main

import (
	"github.com/r3volut1oner/go-karbo/cmd"
	"runtime"
)

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	cmd.Execute()
}
