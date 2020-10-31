package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sethpyle376/cs-statman/statman/grpcserver"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		println(".env not found, assuming non-local environment")
	}

	gs, err := grpcserver.New("4000")

	if err != nil {
		println("ERROR: " + err.Error())
	}

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
