package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type IndexController struct{}

func NewIndexController() Controller {
	return &IndexController{}
}

func (c *IndexController) Mappings() []Mapping {
	mapping := Mapping{Path: "/", Method: "GET", Handler: c.Index()}
	return append([]Mapping{}, mapping)
}

func (c *IndexController) Index() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "index", nil)
	}
}
