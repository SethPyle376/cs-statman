package matchprocessor

import (
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"os"
)

type MatchProcessor struct{}

func New() (*MatchProcessor, error) {
	mp := &MatchProcessor{}
	return mp, nil
}

func (mp *MatchProcessor) ProcessMatch(file *os.File) string {
	f, err := os.Open(file.Name())

	if err != nil {
		println("error", err)
	} else {
		println("good")
	}

	parser := dem.NewParser(f)
	defer parser.Close()

	var message string

	parser.RegisterEventHandler(func(bd events.BombDefused) {
		var bombsite string
		if bd.Site == events.BombsiteA {
			bombsite = "A"
		} else {
			bombsite = "B"
		}
		message += ("Player: " + bd.Player.Name + " defused bomb at site: " + bombsite + "\n")
	})

	err = parser.ParseToEnd()

	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	return message
}
