package main

import (
	"html/template"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	tmpl, err := template.New("react.tmpl").ParseFiles("./ui/html/react.tmpl")
	if err != nil {
		return nil, err
	}

	cache["react.tmpl"] = tmpl

	return cache, nil
}
