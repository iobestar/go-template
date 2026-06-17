package cmd

import (
	"fmt"
	"strings"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name = "template"
)

var (
	Version = "N/A"
	Build   = "N/A"
)

type commander interface {
	getCommand() *cobra.Command
}

type baseCmd struct {
	cmd *cobra.Command
}

func newBaseCmd(cmd *cobra.Command) *baseCmd {
	return &baseCmd{
		cmd: cmd,
	}
}

func (bc *baseCmd) getCommand() *cobra.Command {
	return bc.cmd
}

func (bc *baseCmd) addCommands(commanders ...commander) {
	for _, command := range commanders {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		bc.getCommand().AddCommand(cmd)
	}
}

type rootCmd struct {
	*baseCmd
}

func newRootCmd() *rootCmd {
	cmd := &rootCmd{
		baseCmd: newBaseCmd(&cobra.Command{
			Use:               name,
			Version:           fmt.Sprintf("%s (build: %s)", Version, Build),
			CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		}),
	}

	cmd.addCommands(newConfigCmd())
	cmd.addCommands(newServerCmd())
	return cmd
}

func Execute() error {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
		viper.SetEnvPrefix(name)
		viper.SetConfigType("yaml")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		configFile := viper.GetString("config.file")
		if len(configFile) > 0 {
			viper.SetConfigFile(configFile)
		} else {
			viper.SetConfigName(name + ".yml")
			viper.AddConfigPath(".")
		}

		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileNotFoundError) {
				internal.Fatal(err, 1)
			}
		}
	})
	return newRootCmd().getCommand().Execute()
}
