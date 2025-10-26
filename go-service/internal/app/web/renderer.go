package web

import (
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/iobestar/go-template/go-service/web"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type PageRenderer struct {
	templates map[string]*template.Template
}

func NewPageRenderer(root fs.FS, name string) (*PageRenderer, error) {
	templates, err := parseTemplates(root, name)
	if nil != err {
		return nil, err
	}

	return &PageRenderer{
		templates: templates,
	}, nil
}

func MustNewPageRenderer(root fs.FS, name string) *PageRenderer {
	if v, err := NewPageRenderer(root, name); nil != err {
		panic(err)
	} else {
		return v
	}
}

func parseTemplates(root fs.FS, name string) (map[string]*template.Template, error) {
	var templates []string
	var templateFiles = make(map[string][]string)

	err := fs.WalkDir(root, name, func(path string, d fs.DirEntry, err error) error {
		if strings.EqualFold(path, name) {
			return nil
		}
		if d.IsDir() {
			templates = append(templates, path)
		} else {
			p := "/" + filepath.Dir(path)
			if v, ok := templateFiles[p]; ok {
				templateFiles[p] = append(v, filepath.Base(path))
			} else {
				templateFiles[p] = []string{filepath.Base(path)}
			}
		}
		return nil
	})
	if nil != err {
		return nil, err
	}

	var result = make(map[string]*template.Template)
	for _, t := range templates {
		tn := strings.TrimPrefix(t, name)
		elements := strings.Split(t, "/")

		var patterns []string
		var p string
		for _, e := range elements {
			p = p + "/" + e

			for _, v := range templateFiles[p] {
				patterns = append(patterns, (p + "/" + v)[1:])
			}
		}

		tmpl, err := template.ParseFS(root, patterns...)
		if nil != err {
			return nil, errors.Wrapf(err, "error parsing template: %s", tn)
		}

		result[tn] = tmpl
	}

	return result, nil
}

func (r *PageRenderer) Render(writer io.Writer, name string, data interface{}, _ echo.Context) error {
	if t, ok := r.templates[name]; ok {
		return t.ExecuteTemplate(writer, "base", data)
	}
	return errors.New("template not found")
}

func (r *PageRenderer) Parse(name string, patterns ...string) error {
	if _, ok := r.templates[name]; ok {
		return errors.Errorf("template already exists: %s", name)
	}

	t, err := template.ParseFS(web.Pages, patterns...)
	if err != nil {
		return errors.Wrap(err, "error parsing template")
	}

	r.templates[name] = t
	return nil
}
