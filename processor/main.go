package main

import (
	// "fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sethpyle376/cs-statman/processor/uploadserver"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		println(".env not found, assuming non-local environment")
	}

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
