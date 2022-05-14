package utils

import (
	"html/template"
	"path/filepath"
)

func TemplateParsing(dir string) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		templates[name] = ts
	}
	return templates, nil
}
