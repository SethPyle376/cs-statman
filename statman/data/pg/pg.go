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
		INSERT INTO match_player (matchID, userID, hltv, kills, deaths, adr)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	for _, player := range players {
		_, err := ps.db.Exec(insertStatement, matchID, player.GetSteamID(), player.GetHltv(), player.GetKills(), player.GetDeaths(), player.GetAdr())
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

func (ps *PostgresStore) GetPlayerMatches(playerID int64) ([]int64, error) {
	statement := `
		select matchID from match_player where userID=$1;
	`

	rows, err := ps.db.Query(statement, playerID)
	if err != nil {
		return nil, nil
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
		fmt.Printf("MATCH SET: %d\n", matchID)
	}

	return matchIDs, nil
}
