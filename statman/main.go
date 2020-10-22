package main

import (
	"github.com/sethpyle376/cs-statman/statman/grpcserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gs, err := grpcserver.New("4000")

	// go until told to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)

	select {
	case <-sigs:
	case <-errChan:
		println(err.Error())
	}

	gs.Stop()

}
