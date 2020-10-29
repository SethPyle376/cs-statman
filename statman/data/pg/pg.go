package pg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"os"
	"strconv"
	"time"
)

type PostgresStore struct {
	db *sql.DB
}

func New() (*PostgresStore, error) {
	host := os.Getenv("PG_HOST")
	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DB_NAME")
	sslMode := os.Getenv("PG_SSL_MODE")

	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	instance := PostgresStore{}
	instance.db = db

	return &instance, nil
}

func (ps *PostgresStore) saveMatchData(match *csproto.MatchData) error {
	insertStatement := `
		INSERT INTO match (id, map, date, roundCount)
		VALUES ($1, $2, $3, $4);
	`

	_, err := ps.db.Exec(insertStatement, match.GetMatchID(), match.GetMap(), time.Now(), match.GetRoundCount())

	return err
}

func (ps *PostgresStore) savePlayers(players []*csproto.PlayerData) error {
	insertStatement := `
		INSERT INTO player (id, name)
		VALUES ($1, $2)
		ON CONFLICT(id) DO UPDATE
		SET name=EXCLUDED.name;
	`
	for _, player := range players {
		_, err := ps.db.Exec(insertStatement, player.GetSteamID(), player.GetName())
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *PostgresStore) savePlayerMatchData(matchID int64, players []*csproto.PlayerData) error {
	insertStatement := `
		INSERT INTO match_player (matchID, userID, name, hltv, kills, deaths, adr)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	for _, player := range players {
		_, err := ps.db.Exec(insertStatement, matchID, player.GetSteamID(), player.GetName(), player.GetHltv(), player.GetKills(), player.GetDeaths(), player.GetAdr())
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PostgresStore) SaveMatch(match *csproto.MatchInfo) error {
	err := ps.saveMatchData(match.GetMatchData())
	if err != nil {
		return err
	}

	err = ps.savePlayers(match.GetPlayerData())
	if err != nil {
		return err
	}

	err = ps.savePlayerMatchData(match.GetMatchData().GetMatchID(), match.GetPlayerData())
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStore) GetMatch(matchID int64) (*csproto.MatchInfo, error) {
	matchInfo := &csproto.MatchInfo{}
	matchData := &csproto.MatchData{}

	matchStatement := `
	SELECT map, roundCount FROM match WHERE id=$1`

	row := ps.db.QueryRow(matchStatement, matchID)

	switch err := row.Scan(&matchData.Map, &matchData.RoundCount); err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		break
	default:
		return nil, err
	}

	matchInfo.MatchData = matchData

	playerStatement := `
	SELECT userID, name, hltv, kills, deaths, adr FROM match_player WHERE matchID=$1;`

	playerRows, err := ps.db.Query(playerStatement, matchID)
	if err != nil {
		return nil, err
	}
	defer playerRows.Close()

	for playerRows.Next() {
		playerData := &csproto.PlayerData{}

		playerRows.Scan(&playerData.SteamID, &playerData.Name, &playerData.Hltv, &playerData.Kills, &playerData.Deaths, &playerData.Adr)
		matchInfo.PlayerData = append(matchInfo.PlayerData, playerData)
	}
	return matchInfo, nil
}

func (ps *PostgresStore) GetPlayerMatches(playerID int64) ([]int64, error) {
	statement := `
		select matchID from match_player where userID=$1;
	`

	rows, err := ps.db.Query(statement, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchIDs []int64

	for rows.Next() {
		var matchID int64
		err = rows.Scan(&matchID)
		if err != nil {
			return nil, err
		}
		matchIDs = append(matchIDs, matchID)
	}

	return matchIDs, nil
}

func (ps *PostgresStore) GetPlayerMatchData(playerID int64) ([]*csproto.PlayerMatchData, error) {
	statement := `
		select md.id AS matchID, md.map AS map, md.roundCount AS roundCount,
				mp.userID AS steamID, mp.name AS name, mp.hltv, mp.kills, mp.deaths,
					mp.adr
		FROM match md 
		INNER JOIN match_player mp ON md.id=mp.matchID
		WHERE mp.userID = $1;
	`

	rows, err := ps.db.Query(statement, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := []*csproto.PlayerMatchData{}

	for rows.Next() {
		playerMatchData := &csproto.PlayerMatchData{}
		matchData := &csproto.MatchData{}
		playerData := &csproto.PlayerData{}

		err = rows.Scan(&matchData.MatchID, &matchData.Map, &matchData.RoundCount,
			&playerData.SteamID, &playerData.Name, &playerData.Hltv,
			&playerData.Kills, &playerData.Deaths, &playerData.Adr)

		if err != nil {
			return nil, err
		}
		playerMatchData.MatchData = matchData
		playerMatchData.PlayerData = playerData

		response = append(response, playerMatchData)
	}

	return response, nil
}
