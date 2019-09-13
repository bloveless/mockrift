package main

import "html/template"

type Flash struct {
	Type    string
	Message template.HTML
}

type SharedTemplateData struct {
	AppName     string
	CurrentYear int
	Flash       []Flash
}
