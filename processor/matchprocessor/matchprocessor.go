package matchprocessor

import (
	"context"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

type MatchProcessor struct {
	client csproto.StatmanClient
}

func New() (*MatchProcessor, error) {
	mp := &MatchProcessor{}
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	client := csproto.NewStatmanClient(conn)
	mp.client = client
	return mp, err
}

func (mp *MatchProcessor) ProcessMatch(file *os.File) string {
	f, err := os.Open(file.Name())

	if err != nil {
		println("error", err)
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

	parser.RegisterEventHandler(func(ph events.PlayerHurt) {
		player := ph.Player.Name
		playerId := strconv.FormatUint(ph.Player.SteamID64, 10)

		var attacker string

		if ph.Attacker != nil {
			attacker = ph.Attacker.Name
		} else {
			attacker = "WORLD"
		}

		weapon := ph.Weapon.Type.String()
		damage := ph.HealthDamage
		isLive := parser.GameState().IsMatchStarted()

		if isLive {
			message += (attacker + " damaged " + player + "(ID: " + playerId + ")" + " for " + strconv.Itoa(damage) + " with " + weapon +
				" on round " + strconv.Itoa(parser.GameState().TotalRoundsPlayed()+1) + "\n")
		}
	})

	err = parser.ParseToEnd()

	testMessage := &csproto.MatchInfo{}
	testMatchData := &csproto.MatchData{}
	testMatchData.Map = "inferno"
	testMessage.MatchData = testMatchData

	request := &csproto.SaveMatchRequest{}
	request.MatchInfo = testMessage

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mp.client.SaveMatch(ctx, request)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	return message
}
