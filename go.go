package template

import (
	"context"
	"embed"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"strings"
)

//go:embed go-service
var service embed.FS

func CreateGolangService(name, path string, module string) error {
	if err := os.MkdirAll(path, 0755); nil != err {
		return err
	}

	log.Info("service path created: ", path)

	if err := os.Chdir(path); nil != err {
		return errors.Wrap(err, "unable to change directory")
	}

	code, _, _, err := Execute(context.Background(), Command{
		Name: "go",
		Args: []string{"mod", "init", module},
	})
	if nil != err {
		return errors.Wrap(err, "unable to initialize module")
	}

	if code != 0 {
		return errors.New("unable to initialize module")
	}

	log.Info("module initialized: ", module)

	err = fs.WalkDir(service, "go-service", func(sp string, d fs.DirEntry, err error) error {
		if nil != err {
			return err
		}

		if d.IsDir() && d.Name() == "go-service" {
			return nil
		}

		tp := strings.ReplaceAll(sp, "go-service", path)
		if strings.HasPrefix(sp, "go-service/cmd") {
			tp = strings.ReplaceAll(tp, "template", name)
		}

		if d.IsDir() {
			return os.MkdirAll(tp, 0755)
		}

		f, err := service.Open(sp)
		if nil != err {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		sb, err := io.ReadAll(f)
		if nil != err {
			return err
		}

		tb := strings.ReplaceAll(string(sb), "github.com/iobestar/go-template/go-service", module)
		if "go-service/cmd/template.go" == sp {
			tb = strings.ReplaceAll(tb, "\"template\"", "\""+name+"\"")
		}
		return os.WriteFile(tp, []byte(tb), 0644)
	})
	if nil != err {
		return err
	}

	return addDependencies("github.com/spf13/viper",
		"github.com/fsnotify/fsnotify",
		"github.com/golang-migrate/migrate/v4",
		"github.com/golang-migrate/migrate/v4/database/postgres",
		"github.com/golang-migrate/migrate/v4/database/sqlite3",
		"github.com/golang-migrate/migrate/v4/source",
		"github.com/jmoiron/sqlx",
		"github.com/lib/pq",
		"github.com/mattn/go-sqlite3",
		"github.com/pkg/errors",
		"github.com/sirupsen/logrus",
		"github.com/labstack/echo/v4",
		"github.com/natefinch/lumberjack",
		"github.com/prometheus/client_golang/prometheus",
		"github.com/prometheus/client_golang/prometheus/promauto",
		"github.com/spf13/cobra",
		"github.com/mitchellh/go-homedir",
		"gopkg.in/yaml.v3",
		"github.com/mattn/go-tty")
}

func addDependencies(modules ...string) error {
	for _, m := range modules {
		code, _, _, err := Execute(context.Background(), Command{Name: "go", Args: []string{"get", m}})
		if nil != err {
			return errors.Wrap(err, fmt.Sprintf("unable to get dependency: %s", m))
		}

		if code != 0 {
			return errors.New(fmt.Sprintf("unable to get dependency: %s", m))
		}
		log.Info("dependency added: ", m)
	}
	return nil
}
