package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

func handleHome(ctx *gin.Context) {
	d := &homeData{
		Shared: getSharedTemplateData(),
		Body:   "",
	}

	ctx.HTML(http.StatusOK, "home.page.tmpl", d)
}
