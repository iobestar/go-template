package app

import (
	"github.com/iobestar/go-template/go-service/internal/app/web"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	config Config

	dbx         *sqlx.DB
	controllers []web.Controller
}

func NewContext(config Config, dbx *sqlx.DB) Context {
	return Context{
		config:      config,
		dbx:         dbx,
		controllers: make([]web.Controller, 0),
	}
}

func (c *Context) register(any any) {
	if value, ok := any.(web.Controller); ok {
		c.controllers = append(c.controllers, value)
	}
}

func (c *Context) Init() error {
	c.register(web.NewIndexController())
	c.register(web.NewAboutController())
	return nil
}

func (c *Context) Controllers() []web.Controller {
	result := make([]web.Controller, 0)
	for _, controller := range c.controllers {
		result = append(result, controller)
	}
	return result
}

func (c *Context) Start() error {
	return nil
}

func (c *Context) Stop() error {
	return nil
}
