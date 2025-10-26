package database

import (
	"fmt"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.name", name)
	viper.SetDefault("database.postgres.user", name)
	viper.SetDefault("database.postgres.password", name)
	viper.SetDefault("database.postgres.params", make(map[string]string))
	viper.SetDefault("database.postgres.max_idle", 10)
	viper.SetDefault("database.postgres.max_open", 50)
}

type PostgresConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Params   map[string]string
	MaxIdle  int
	MaxOpen  int
}

func NewPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetInt("database.postgres.port"),
		Name:     viper.GetString("database.postgres.name"),
		User:     viper.GetString("database.postgres.user"),
		Password: viper.GetString("database.postgres.password"),
		Params:   viper.GetStringMapString("database.postgres.params"),
		MaxIdle:  viper.GetInt("database.postgres.max_idle"),
		MaxOpen:  viper.GetInt("database.postgres.max_open"),
	}
}

func (c PostgresConfig) Url() string {
	params := make([]string, 0)
	for k, v := range c.Params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	if len(params) == 0 {
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Name)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", c.User, c.Password, c.Host, c.Port, c.Name, strings.Join(params, "&"))
}

func ConnectPostgres(config PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", config.Url())
	if nil != err {
		return nil, err
	}
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetMaxOpenConns(config.MaxOpen)

	if err := db.Ping(); nil != err {
		defer func() {
			_ = db.Close()
		}()
		return nil, err
	}
	return db, nil
}
