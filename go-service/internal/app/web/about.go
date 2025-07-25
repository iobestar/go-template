package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type AboutController struct{}

func NewAboutController() Controller {
	return &AboutController{}
}

func (c *AboutController) Mappings() []Mapping {
	mapping := Mapping{Path: "/about", Method: "GET", Handler: c.About()}
	return append([]Mapping{}, mapping)
}

func (c *AboutController) About() echo.HandlerFunc {
	return func(ctx echo.Context) error {

		return ctx.Render(http.StatusOK, "about", nil)
	}
}
