package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/mitchellh/go-homedir"
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
		if config := os.Getenv(strings.ToUpper(name) + "_CONFIG"); len(config) > 0 {
			viper.SetConfigFile(config)
		} else {
			homeDir, err := homedir.Dir()
			if err != nil {
				internal.FatalError(errors.Wrap(err, "unable to resolve home dir"))
			}
			viper.SetConfigName(name + ".yml")
			viper.SetConfigType("yaml")

			viper.AddConfigPath(".")
			viper.AddConfigPath(path.Join(homeDir, "."+name))
		}
		viper.AutomaticEnv()
		viper.SetEnvPrefix(name)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileNotFoundError) {
				internal.FatalError(errors.Wrap(err, "unable to read config"))
			}
		}
	})
	return newRootCmd().getCommand().Execute()
}
