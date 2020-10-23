package main

import (
	// "fmt"

	"github.com/sethpyle376/cs-statman/processor/uploadserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	us, err := uploadserver.New()

	go us.Start()

	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-sigs:
	}

	println("shutting down")
}
