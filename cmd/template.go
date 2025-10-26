package cmd

import (
	"fmt"

	"github.com/iobestar/go-template"
	"github.com/spf13/cobra"
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

func (c *baseCmd) getCommand() *cobra.Command {
	return c.cmd
}

func (c *baseCmd) addCommands(commanders ...commander) {
	for _, command := range commanders {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		c.getCommand().AddCommand(cmd)
	}
}

func newRootCmd() *baseCmd {
	cmd := newBaseCmd(&cobra.Command{
		Use:               "template",
		Version:           fmt.Sprintf("%s (build: %s)", "N/A", "N/A"),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	})
	cmd.addCommands(newCreateCmd())
	return cmd
}

func Execute() error {
	cobra.OnInitialize(func() {
		template.SetupLogging()
	})
	return newRootCmd().getCommand().Execute()
}
