package server

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("server.network", "tcp")
	viper.SetDefault("server.address", ":3200")
	viper.SetDefault("server.key_path", "")
	viper.SetDefault("server.cert_path", "")
}

type Config struct {
	Network  string
	Address  string
	KeyPath  string
	CertPath string
}

func NewConfig() Config {
	return Config{
		Network:  viper.GetString("server.network"),
		Address:  viper.GetString("server.address"),
		KeyPath:  viper.GetString("server.key_path"),
		CertPath: viper.GetString("server.cert_path"),
	}
}
