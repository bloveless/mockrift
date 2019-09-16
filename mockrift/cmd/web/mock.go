package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mockrift/pkg/helper"
	"net/http"
)

func addMockingRoutes(router *gin.Engine, client *http.Client, recordOnly bool) {
	var a *App
	router.Any("/m/:app/*path", func(ctx *gin.Context) {
		appName := ctx.Param("app")
		path := ctx.Param("path")

		reqBody, reqBodyErr := ioutil.ReadAll(ctx.Request.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		if a.Name != appName {
			a = getAppFromFile(appName)
		}

		storedRes := a.findResponseByRequestParams(ctx.Request.Method, path, reqBody)
		if !recordOnly && storedRes != nil {
			fmt.Println("Sending response from memory")
			helper.CopyHeaders(storedRes.Header, ctx)
			ctx.Writer.WriteHeader(storedRes.StatusCode)
			helper.CopyBodyToClient(ctx.Writer, storedRes.Body)
		} else {
			fmt.Println("Proxying response to real backend")
			cReq := helper.CreateClientRequest(ctx.Request.Method, "http://host.docker.internal:4000"+path, reqBody)
			defer cReq.Body.Close()

			cRes := helper.DoClientRequest(client, cReq)

			cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
			if cBodyErr != nil {
				log.Fatal(cBodyErr)
			}

			fmt.Printf("Client response body: %v\n", string(cBody))

			helper.CopyHeaders(cReq.Header, ctx)
			ctx.Writer.WriteHeader(cRes.StatusCode)
			helper.CopyBodyToClient(ctx.Writer, cBody)

			a.storeResponseAndRequest(ctx.Request, path, reqBody, cRes, cBody)
		}
	})
}
