package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Pages struct {
	templateCache map[string]*template.Template
}

func getSharedTemplateData() *SharedTemplateData {
	return &SharedTemplateData{
		AppName:     "Mockrift",
		CurrentYear: time.Now().Year(),
		Flash:       "",
	}
}

type homeData struct {
	Shared *SharedTemplateData
	Body   string
}

func (p *Pages) handleHome(w http.ResponseWriter, req *http.Request) {
	t, ok := p.templateCache["home.page.tmpl"]
	if !ok {
		log.Fatal("Unable to find home page template")
	}

	fmt.Println("Handling Home")

	buf := new(bytes.Buffer)

	d := &homeData{
		Shared: getSharedTemplateData(),
		Body:   "",
	}

	err := t.Execute(buf, d)
	if err != nil {
		log.Fatal("Unable to execute template: " + err.Error())
	}

	buf.WriteTo(w)
}
