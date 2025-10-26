package web

import (
	"io/fs"
	"net/http"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/iobestar/go-template/go-service/web"
	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, appCtx internal.AppContext) error {
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", getStaticHandler())))

	e.Renderer = MustNewPageRenderer(web.Pages, "pages")
	return configure(e, appCtx)
}

func getStaticHandler() http.Handler {
	staticSub, err := fs.Sub(web.Static, "static")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(staticSub))
}

func configure(e *echo.Echo, appCtx internal.AppContext) error {
	controllers := internal.GetAllByType[Controller](appCtx)

	for _, controller := range controllers {
		for _, route := range controller.Routes() {
			e.Match(route.Methods, route.Path, route.Handler)
		}
	}
	return nil
}
