package main

import (
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"github.com/sethpyle376/cs-statman/uploadserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	us, err := uploadserver.New()
	f, err := os.Open("test/test.dem")

	if err != nil || us == nil {
		println("failed to open file")
		panic(err)
	}

	defer f.Close()

	go us.Start()

	p := dem.NewParser(f)
	defer p.Close()

	p.RegisterEventHandler(func(e events.Kill) {
		var hs string
		if e.IsHeadshot {
			hs = " (HS)"
		}
		var wallbang string
		if e.IsWallBang() {
			wallbang = " (WB)"
		}
		fmt.Printf("%s <%v%s%s> %s\n", e.Killer, e.Weapon.String(), hs, wallbang, e.Victim)
	})
	err = p.ParseToEnd()
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
