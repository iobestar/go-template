package web

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	Methods []string
	Path    string
	Handler echo.HandlerFunc
}

type Controller interface {
	Routes() []Route
}
