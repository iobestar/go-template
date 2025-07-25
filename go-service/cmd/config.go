package cmd

import (
	"github.com/iobestar/go-template/go-service/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Config string
	Format string
}

func NewConfig() Config {
	return Config{
		Config: viper.GetString("config"),
	}
}

func newConfigCmd() *baseCmd {
	return newBaseCmd(&cobra.Command{
		Use:   "config",
		Short: "Show effective configuration",
		RunE: func(_ *cobra.Command, args []string) error {
			config := make(map[string]interface{})
			if err := viper.Unmarshal(&config); nil != err {
				return err
			}

			internal.WriteFormat(config)
			return nil
		},
	})
}
