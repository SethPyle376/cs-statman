package matchprocessor

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"io"
	"os"
	"time"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
)

type MatchProcessor struct {
	client csproto.StatmanClient
}

func New() (*MatchProcessor, error) {
	mp := &MatchProcessor{}

	var grpcConn *grpc.ClientConn

	cert, ok := os.LookupEnv("TLS_CERT_LOCATION")

	var err error

	host, hostOk := os.LookupEnv("STATMAN_HOST")

	if !hostOk {
		host = "localhost:4000"
	}

	if !ok {
		grpcConn, err = grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		grpcOpt, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return nil, err
		}
		grpcConn, err = grpc.Dial(host, grpc.WithTransportCredentials(grpcOpt))
	}

	client := csproto.NewStatmanClient(grpcConn)
	mp.client = client
	return mp, nil
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

	var playerTeams map[uint64]common.Team
	playerTeams = make(map[uint64]common.Team)

	var playerKills map[uint64]int
	playerKills = make(map[uint64]int)

	var playerDeaths map[uint64]int
	playerDeaths = make(map[uint64]int)

	var roundKills map[int][]*csproto.Kill
	roundKills = make(map[int][]*csproto.Kill)

	var roundWinner map[int]common.Team
	roundWinner = make(map[int]common.Team)

	currentRound := 0

	parser.RegisterEventHandler(func(ph events.PlayerHurt) {
		playerTeams[ph.Player.SteamID64] = ph.Player.Team
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
				}
			}
		}
	})

	parser.RegisterEventHandler(func(rs events.RoundStart) {
		currentRound++
	})

	parser.RegisterEventHandler(func(rs events.RoundEnd) {
		round := parser.GameState().TotalRoundsPlayed()
		roundWinner[round+1] = rs.Winner
	})

	parser.RegisterEventHandler(func(pc events.PlayerConnect) {
		playerNames[pc.Player.SteamID64] = pc.Player.Name
		playerTeams[pc.Player.SteamID64] = pc.Player.Team
	})

	parser.RegisterEventHandler(func(k events.Kill) {
		if parser.GameState().IsMatchStarted() {
			if k.Killer != nil {
				playerKills[k.Killer.SteamID64]++

				kill := &csproto.Kill{}
				kill.KillerID = k.Killer.SteamID64
				kill.VictimID = k.Victim.SteamID64
				roundKills[currentRound] = append(roundKills[currentRound], kill)
			}
			playerDeaths[k.Victim.SteamID64]++
		}
	})

	parser.ParseToEnd()

	totalRounds := parser.GameState().TotalRoundsPlayed()
	matchData.RoundCount = int32(totalRounds)

	var playerData []*csproto.PlayerData

	for k, v := range playerTotalDamage {
		player := &csproto.PlayerData{}
		player.Adr = float32(v) / float32(totalRounds)
		player.Name = playerNames[k]
		player.Team = int32(playerTeams[k])
		player.Kills = int32(playerKills[k])
		player.Deaths = int32(playerDeaths[k])
		player.SteamID = int64(k)
		playerData = append(playerData, player)
	}

	matchInfo.RoundData = make([]*csproto.RoundData, totalRounds+1)

	for i := 0; i <= totalRounds; i++ {
		roundData := &csproto.RoundData{}
		roundData.WinningTeam = int32(roundWinner[i])

		for _, kill := range roundKills[i] {
			roundData.Kills = append(roundData.Kills, kill)
		}
		matchInfo.RoundData[i] = roundData
	}

	matchInfo.MatchData = matchData
	matchInfo.PlayerData = playerData

	return matchInfo, err
}

func (mp *MatchProcessor) ProcessMatch(file *os.File) error {
	h := md5.New()

	file, err := os.Open(file.Name())

	if _, err := io.Copy(h, file); err != nil {
		return err
	}

	hash := h.Sum(nil)

	matchHash := binary.BigEndian.Uint64(hash)

	matchInfo, err := mp.getMatchStats(file)

	if err != nil {
		return err
	}

	request := &csproto.SaveMatchRequest{}
	request.MatchInfo = matchInfo
	request.MatchInfo.MatchData.MatchID = int64(matchHash)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mp.client.SaveMatch(ctx, request)

	defer file.Close()
	defer os.Remove(file.Name())

	return nil
}
