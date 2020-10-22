package matchprocessor

import (
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"github.com/sethpyle376/cs-statman/csproto"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

type MatchProcessor struct {
	conn *grpc.ClientConn
}

func New() (*MatchProcessor, error) {
	mp := &MatchProcessor{}
	var opts []grpc.DialOption
	conn, err := grpc.Dial("localhost:4000", opts...)
	mp.conn = conn
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

	if err != nil {
		panic(err)
	}

	mp.conn.Se

	defer file.Close()
	defer os.Remove(file.Name())

	return message
}
