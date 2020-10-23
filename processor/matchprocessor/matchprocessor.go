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

func (mp *MatchProcessor) getMatchStats(file *os.File) (string, error) {
	f, err := os.Open(file.Name())

	parser := dem.NewParser(f)

	var message string

	parser.RegisterEventHandler(func(ph events.PlayerHurt) {
		player := ph.Player.Name
		playerID := strconv.FormatUint(ph.Player.SteamID64, 10)

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
			message += (attacker + " damaged " + player + "(ID: " + playerID + ")" + " for " + strconv.Itoa(damage) + " with " + weapon +
				" on round " + strconv.Itoa(parser.GameState().TotalRoundsPlayed()+1) + "\n")
		}
	})

	parser.ParseToEnd()
	return message, err
}

func (mp *MatchProcessor) ProcessMatch(file *os.File) string {
	message, err := mp.getMatchStats(file)

	if err != nil {
		println(err.Error())
		return "error\n"
	}

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
