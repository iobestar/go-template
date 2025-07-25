package logging

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.file.enabled", false)
	viper.SetDefault("logging.file.name", "")
	viper.SetDefault("logging.file.max_size", "10")
	viper.SetDefault("logging.file.max_age", "10")
	viper.SetDefault("logging.file.max_backups", "10")
}

type Config struct {
	Level string
	File  FileConfig
}

type FileConfig struct {
	Enabled    bool
	Name       string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

func NewConfig() Config {
	return Config{
		Level: viper.GetString("logging.level"),
		File: FileConfig{
			Enabled:    viper.GetBool("logging.file.enabled"),
			Name:       viper.GetString("logging.file.name"),
			MaxSize:    viper.GetInt("logging.file.max_size"),
			MaxAge:     viper.GetInt("logging.file.max_age"),
			MaxBackups: viper.GetInt("logging.file.max_backups"),
		},
	}
}
