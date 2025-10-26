package controllers

import (
	"net/http"

	"github.com/iobestar/go-template/go-service/internal/app/web"
	"github.com/labstack/echo/v4"
)

type IndexController struct {
}

func NewIndexController() web.Controller {
	return &IndexController{}
}

func (c *IndexController) Routes() []web.Route {
	return append([]web.Route{},
		web.Route{Methods: []string{"GET"}, Path: "/", Handler: c.Index()},
	)
}

func (c *IndexController) Index() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "/index", nil)
	}
}
