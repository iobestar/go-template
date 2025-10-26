package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func Connect(config Config) (*sqlx.DB, error) {
	switch c := config.(type) {
	case SqliteConfig:
		return ConnectSqlite(c)
	case PostgresConfig:
		return ConnectPostgres(c)
	}
	return nil, errors.New(fmt.Sprintf("unknown config type: %T", config))
}
