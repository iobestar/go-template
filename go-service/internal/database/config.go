package database

import (
	"fmt"
	"github.com/spf13/viper"
)

const name = "template"

func init() {
	viper.SetDefault("database.type", "none")
}

func getType() string {
	return viper.GetString("database.type")
}

type Config interface {
	Url() string
}

func NewConfig() Config {
	t := getType()
	switch t {
	case "sqlite":
		return NewSqliteConfig()
	case "postgres":
		return NewPostgresConfig()
	case "none":
		return NewNoneConfig()
	}
	panic(fmt.Sprintf("unknown database type: %s", t))
}
