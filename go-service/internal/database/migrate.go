package database

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type EmbedMigrationDriver struct {
	fs         embed.FS
	path       string
	migrations *source.Migrations
}

func NewEmbedMigrationDriver(fs embed.FS, path string) (source.Driver, error) {
	files, err := fs.ReadDir(path)
	if nil != err {
		return nil, err
	}

	migrations := source.NewMigrations()
	for _, file := range files {
		if migration, err := source.DefaultParse(file.Name()); nil != err {
			continue
		} else {
			if !migrations.Append(migration) {
				return nil, source.ErrDuplicateMigration{
					Migration: *migration,
				}
			}
		}
	}

	return &EmbedMigrationDriver{
		fs:         fs,
		path:       path,
		migrations: migrations,
	}, nil
}

func (e *EmbedMigrationDriver) First() (version uint, err error) {
	if v, ok := e.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: e.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *EmbedMigrationDriver) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := e.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: e.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *EmbedMigrationDriver) Next(version uint) (nextVersion uint, err error) {
	if v, ok := e.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: e.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *EmbedMigrationDriver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Up(version); ok {
		body, err := e.fs.Open(filepath.Join(e.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(body), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.path, Err: os.ErrNotExist}
}

func (e *EmbedMigrationDriver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Down(version); ok {
		body, err := e.fs.Open(filepath.Join(e.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(body), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.path, Err: os.ErrNotExist}
}

func (e *EmbedMigrationDriver) Close() error {
	return nil
}

func (e *EmbedMigrationDriver) Open(url string) (source.Driver, error) {
	panic("not implemented")
}

func Migrate(migration embed.FS, config Config) error {
	sourceDriver, err := NewEmbedMigrationDriver(migration, "migration")
	if nil != err {
		return err
	}

	var m *migrate.Migrate
	if c, ok := config.(SqliteConfig); ok {
		d, err := ConnectSqlite(c)
		if nil != err {
			return err
		}

		databaseDriver, err := sqlite3.WithInstance(d.DB, &sqlite3.Config{})
		if nil != err {
			return err
		}
		m, err = migrate.NewWithInstance("embedded", sourceDriver, "sqlite", databaseDriver)
		if err != nil {
			return err
		}
	} else {
		m, err = migrate.NewWithSourceInstance("embedded", sourceDriver, config.Url())
		if err != nil {
			return err
		}
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if nil != srcErr {
			log.WithError(err).Warn("error closing migration driver")
		}
		if nil != dbErr {
			log.WithError(err).Warn("error closing migration database")
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}
