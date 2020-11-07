package data

import (
	"github.com/sethpyle376/cs-statman/pkg/csproto"
)

type Store interface {
	SaveMatch(match *csproto.MatchInfo) error
	GetPlayerMatches(playerID int64) ([]int64, error)
	GetMatch(matchID int64) (*csproto.MatchInfo, error)
	GetRecentMatches() ([]*csproto.MatchData, error)
	GetPlayerMatchData(playerID int64) ([]*csproto.PlayerMatchData, error)
}
