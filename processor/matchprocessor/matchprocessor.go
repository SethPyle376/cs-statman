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

func (mp *MatchProcessor) getMatchStats(file *os.File) (*csproto.MatchInfo, error) {
	matchInfo := &csproto.MatchInfo{}
	f, err := os.Open(file.Name())

	parser := dem.NewParser(f)

	headers, err := parser.ParseHeader()

	matchData := &csproto.MatchData{}
	matchData.Map = headers.MapName

	var playerTotalDamage map[uint64]int
	playerTotalDamage = make(map[uint64]int)

	var playerNames map[uint64]string
	playerNames = make(map[uint64]string)

	parser.RegisterEventHandler(func(ph events.PlayerHurt) {
		if parser.GameState().IsMatchStarted() && ph.Attacker != nil {
			if ph.Attacker != ph.Player {
				if ph.Attacker.Team != ph.Player.Team {
					var actualDamage int

					if ph.Player.Health() < ph.HealthDamage {
						actualDamage = ph.Player.Health()
					} else {
						actualDamage = ph.HealthDamage
					}

					playerTotalDamage[ph.Attacker.SteamID64] += actualDamage
					println(ph.Attacker.Name + " damaged " + ph.Player.Name + " for " + strconv.Itoa(actualDamage) + " with " + ph.Weapon.String())
				}
			}
		}
	})

	parser.RegisterEventHandler(func(pc events.PlayerConnect) {
		playerNames[pc.Player.SteamID64] = pc.Player.Name
	})

	parser.ParseToEnd()

	totalRounds := parser.GameState().TotalRoundsPlayed()
	matchData.RoundCount = int32(totalRounds)

	println(strconv.Itoa(totalRounds) + " rounds played")

	var playerData []*csproto.PlayerData

	for k, v := range playerTotalDamage {
		player := &csproto.PlayerData{}
		player.Adr = float32(v) / float32(totalRounds)
		player.Name = playerNames[k]
		playerData = append(playerData, player)
	}

	matchInfo.MatchData = matchData
	matchInfo.PlayerData = playerData

	return matchInfo, err
}

func (mp *MatchProcessor) ProcessMatch(file *os.File) error {
	matchInfo, err := mp.getMatchStats(file)

	if err != nil {
		return err
	}

	request := &csproto.SaveMatchRequest{}
	request.MatchInfo = matchInfo

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mp.client.SaveMatch(ctx, request)

	defer file.Close()
	defer os.Remove(file.Name())

	return nil
}
