package database

import (
	"fmt"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("database.sqlite.path", name+".db")
	viper.SetDefault("database.sqlite.params", getDefaultParams())
}

func getDefaultParams() map[string]string {
	return map[string]string{
		"_timeout": "5000",
		"_journal": "wal",
		"mode":     "rwc",
	}
}

type SqliteConfig struct {
	Path   string
	Params map[string]string
}

func NewSqliteConfig() SqliteConfig {
	return SqliteConfig{
		Path:   viper.GetString("database.sqlite.path"),
		Params: viper.GetStringMapString("database.sqlite.params"),
	}
}

func (c SqliteConfig) Url() string {
	params := make([]string, 0)
	for k, v := range c.Params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	if len(params) == 0 {
		return fmt.Sprintf("file:%s", c.Path)
	}
	return fmt.Sprintf("file:%s?&%s", c.Path, strings.Join(params, "&"))
}

func ConnectSqlite(config SqliteConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", config.Url())
	if nil != err {
		return nil, err
	}

	if err := db.Ping(); nil != err {
		defer func() {
			_ = db.Close()
		}()
		return nil, err
	}
	return db, nil
}
