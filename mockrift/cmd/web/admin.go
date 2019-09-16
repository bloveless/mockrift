package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func populateSharedTemplate(data *SharedTemplateData) {
	data.AppName = "Mockrift"
	data.CurrentYear = time.Now().Year()
}

type homeData struct {
	SharedTemplateData
	Body string
}

func addAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.GET("/", handleHome)
	}
}

func handleHome(ctx *gin.Context) {
	getAllApps()

	d := homeData{
		Body: "",
	}

	populateSharedTemplate(&d.SharedTemplateData)

	d.Flash = append(d.Flash, Flash{
		Type:    "success",
		Message: "This is a success alert&mdash;check it out!",
	})

	d.Flash = append(d.Flash, Flash{
		Type:    "danger",
		Message: "This is a danger alert&mdash;check it out!",
	})

	log.Println("Test again")

	ctx.HTML(http.StatusOK, "home.page.tmpl", d)
}
