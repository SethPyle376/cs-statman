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

func (ps *PostgresStore) SaveMatch(match *csproto.MatchInfo) error {
	return nil
}
