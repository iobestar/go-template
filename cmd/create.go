package cmd

import (
	"github.com/iobestar/go-template"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func newCreateCmd() *baseCmd {
	cmd := newBaseCmd(&cobra.Command{
		Use:   "create",
		Short: "Create from template",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if nil != err {
				return err
			}

			name, err := template.ReadInput("Enter service name", template.InputOptions{Required: true})
			if nil != err {
				return err
			}

			path, err := template.ReadInput("Enter service path", template.InputOptions{Default: name, Required: true})
			if nil != err {
				return err
			}

			module, err := template.ReadInput("Enter module name", template.InputOptions{Default: "github.com/iobestar/" + name, Required: true})
			if nil != err {
				return err
			}

			return template.CreateGolangService(name, filepath.Join(wd, path), module)
		},
	})
	return cmd
}
