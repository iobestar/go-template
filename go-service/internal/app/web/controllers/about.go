package controllers

import (
	"net/http"

	"github.com/iobestar/go-template/go-service/internal/app/web"
	"github.com/labstack/echo/v4"
)

type AboutController struct{}

func NewAboutController() web.Controller {
	return &AboutController{}
}

func (c *AboutController) Routes() []web.Route {
	return append([]web.Route{},
		web.Route{Methods: []string{"GET"}, Path: "/about", Handler: c.About()},
	)
}

func (c *AboutController) About() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, ctx.Path(), nil)
	}
}
