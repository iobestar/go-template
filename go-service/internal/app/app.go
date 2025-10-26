package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/iobestar/go-template/go-service/internal/app/config"
	"github.com/iobestar/go-template/go-service/internal/app/web"
	"github.com/iobestar/go-template/go-service/internal/app/web/controllers"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type app struct {
	echo     *echo.Echo
	appCnt   internal.AppContainer
	disposer internal.Disposer
}

func NewHttp() (internal.HttpApp, error) {
	a, err := New()
	if err != nil {
		return nil, err
	}

	h, ok := a.(internal.HttpApp)
	if !ok {
		return nil, errors.New("app is not type of internal.HttpApp")
	}
	return h, nil
}

func New() (internal.App, error) {
	return &app{
		echo:   echo.New(),
		appCnt: internal.NewAppContainer(),
	}, nil
}

func (a *app) GetContext() internal.AppContext {
	return a.appCnt
}

func (a *app) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.echo.ServeHTTP(writer, request)
}

func (a *app) Run(ctx context.Context, wg *sync.WaitGroup, disposer internal.Disposer) error {
	appCnt := a.appCnt
	appCnt.MustRegister(internal.AppComponent{Name: "config", Value: config.New()})

	if err := a.setupDB(ctx, wg, disposer); nil != err {
		return err
	}

	if err := a.setupApp(ctx, wg, disposer); nil != err {
		return err
	}

	if err := a.setupApi(); nil != err {
		return err
	}

	if err := a.setupWeb(); nil != err {
		return err
	}

	if err := a.appCnt.Start(); nil != err {
		return err
	}
	disposer.Attach(func() {
		if err := a.appCnt.Stop(); nil != err {
			log.WithError(err).Error("error stopping app container")
		}
	})

	if err := a.appCnt.Run(ctx, wg, disposer); nil != err {
		return err
	}

	return skunkworks(a.appCnt)
}

func (a *app) setupDB(_ context.Context, _ *sync.WaitGroup, _ internal.Disposer) error {
	// setup database connections here
	return nil
}

func (a *app) setupApp(_ context.Context, _ *sync.WaitGroup, _ internal.Disposer) error {
	// setup application services here
	return nil
}

func (a *app) setupApi() error {
	// setup api routes here
	return nil
}

func (a *app) setupWeb() error {
	// setup web controllers here
	appCnt := a.appCnt

	appCnt.MustRegister(internal.AppComponent{
		Name:  "indexController",
		Value: controllers.NewIndexController(),
	})

	appCnt.MustRegister(internal.AppComponent{
		Name:  "aboutController",
		Value: controllers.NewAboutController(),
	})

	return web.Setup(a.echo, a.appCnt)
}
