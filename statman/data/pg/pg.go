package pg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"os"
	"strconv"
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
	return nil
}

func (ps *PostgresStore) SaveMatch(match *csproto.MatchInfo) error {
	err := ps.saveMatchData(match.GetMatchData())

	fmt.Printf("MATCH ID: %d\n", match.GetMatchData().GetMatchID())

	for _, element := range match.GetPlayerData() {
		println("PLAYER: " + element.GetName())
		fmt.Printf("KILLS: %d\n", element.Kills)
		fmt.Printf("DEATHS: %d\n", element.Deaths)
		fmt.Printf("ADR: %f\n", element.Adr)
		fmt.Printf("\n")
	}

	if err != nil {
		return err
	}

	return nil
}
