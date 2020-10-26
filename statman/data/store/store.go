package store

import (
	"github.com/sethpyle376/cs-statman/statman/data"
	"github.com/sethpyle376/cs-statman/statman/data/pg"
)

func New(storeType string) (data.Store, error) {
	switch storeType {
	case "postgres":
		{
			ps, err := pg.New()
			return ps, err
		}
	}
	return nil, nil
}
