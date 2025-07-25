package web

import (
	"github.com/iobestar/go-template/go-service/web"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"io/fs"
	"net/http"
)

type Mapping struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Controller interface {
	Mappings() []Mapping
}

type Renderer struct {
	templates map[string]*template.Template
}

func (r *Renderer) Render(writer io.Writer, name string, data interface{}, _ echo.Context) error {
	if t, ok := r.templates[name]; ok {
		return t.ExecuteTemplate(writer, "base", data)
	}
	return errors.New("template not found")
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
	}
}

func (r *Renderer) Parse(name string, patterns ...string) error {
	if _, ok := r.templates[name]; ok {
		return errors.Errorf("template already exists: %s", name)
	}
	t, err := template.ParseFS(web.Templates, patterns...)
	if err != nil {
		return errors.Wrap(err, "error parsing template")
	}
	r.templates[name] = t
	return nil
}

func (r *Renderer) MustParse(name string, patterns ...string) {
	if err := r.Parse(name, patterns...); err != nil {
		panic(err)
	}
}

func staticHandler() http.Handler {
	staticSub, err := fs.Sub(web.Static, "static")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(staticSub))
}

func initRenderer() echo.Renderer {
	renderer := NewRenderer()
	renderer.MustParse("index", "templates/base.gohtml", "templates/index.gohtml")
	renderer.MustParse("about", "templates/base.gohtml", "templates/about.gohtml")
	return renderer
}

func Setup(e *echo.Echo, controllers ...Controller) error {
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticHandler())))
	e.Renderer = initRenderer()

	for _, controller := range controllers {
		for _, mapping := range controller.Mappings() {
			e.Add(mapping.Method, mapping.Path, mapping.Handler, mapping.Middlewares...)
		}
	}
	return nil
}
