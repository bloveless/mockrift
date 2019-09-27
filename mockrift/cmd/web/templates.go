package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
)

type templateData struct {
	ReactManifest map[string]string
}

func addDefaultData(td *templateData) {
	buf, err := ioutil.ReadFile("./ui/static/react/manifest.json")
	if err != nil {
		log.Fatal("Unable to read manifest.json: " + err.Error())
	}

	uErr := json.Unmarshal(buf, &td.ReactManifest)
	if uErr != nil {
		log.Fatal("Unable to parse manifest.json: " + uErr.Error())
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	tmpl, err := template.New("react.tmpl").ParseFiles("./ui/html/react.tmpl")
	if err != nil {
		return nil, err
	}

	cache["react.tmpl"] = tmpl

	return cache, nil
}
