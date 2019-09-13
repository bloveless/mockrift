package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func getSharedTemplate() SharedTemplateData {
	return SharedTemplateData{
		AppName:     "Mockrift",
		CurrentYear: time.Now().Year(),
	}
}

type homeData struct {
	Shared SharedTemplateData
	Body   string
}

func handleHome(ctx *gin.Context) {
	d := homeData{
		Shared: getSharedTemplate(),
		Body:   "",
	}

	d.Shared.Flash = append(d.Shared.Flash, Flash{
		Type:    "success",
		Message: "This is a success alert&mdash;check it out!",
	})

	d.Shared.Flash = append(d.Shared.Flash, Flash{
		Type:    "danger",
		Message: "This is a danger alert&mdash;check it out!",
	})

	log.Println("Test again")

	ctx.HTML(http.StatusOK, "home.page.tmpl", d)
}
