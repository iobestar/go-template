package app

import (
	"context"
	"github.com/iobestar/go-template/go-service/internal"
	"github.com/iobestar/go-template/go-service/internal/app/db"
	"github.com/iobestar/go-template/go-service/internal/app/web"
	"github.com/iobestar/go-template/go-service/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type App interface {
	internal.Runner
	http.Handler
}

type app struct {
	echo    *echo.Echo
	context Context

	disposer internal.Disposer
}

func New() (App, error) {
	disposer := internal.NewDisposer()

	databaseCfg := database.NewConfig()
	if _, ok := databaseCfg.(database.NoneConfig); !ok {
		if err := database.Migrate(db.Migration, databaseCfg); nil != err {
			return nil, errors.Wrap(err, "error database migrate")
		}

		dbx, err := database.Connect(databaseCfg)
		if nil != err {
			return nil, errors.Wrap(err, "error database connect")
		}

		disposer.Attach(func() {
			if err := dbx.Close(); nil != err {
				log.WithError(err).Error("error database close")
			}
		})
	}

	e := echo.New()

	appCtx := NewContext(NewConfig(), nil) // Pass the database connection to the context if needed
	if err := appCtx.Init(); nil != err {
		return nil, errors.Wrap(err, "error app context init")
	}

	if err := web.Setup(e, appCtx.Controllers()...); nil != err {
		return nil, errors.Wrap(err, "error setup web controllers")
	}

	if err := appCtx.Start(); nil != err {
		return nil, errors.Wrap(err, "error app context start")
	}
	disposer.Attach(func() {
		if err := appCtx.Stop(); nil != err {
			log.WithError(err).Error("error app context stop")
		}
	})

	return &app{
		echo:     e,
		context:  appCtx,
		disposer: disposer,
	}, nil
}

func (a *app) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.echo.ServeHTTP(writer, request)
}

func (a *app) Run(_ context.Context, _ *sync.WaitGroup, disposer internal.Disposer) error {
	disposer.Attach(func() {
		a.disposer.Dispose()
	})
	return nil
}
