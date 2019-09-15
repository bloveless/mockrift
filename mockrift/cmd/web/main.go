package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mockrift/pkg/helper"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", ":3499", "The address to run the mockrift server")
	flag.Parse()

	// urlRegex := regexp.MustCompile("/(.+?)/(.*)")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	recordOnly := false

	router := gin.Default()
	router.Use(cors.Default())
	router.LoadHTMLGlob("./ui/html/*")

	router.StaticFS("/static", http.Dir("./ui/static"))
	router.GET("/admin", handleHome)

	router.Any("/m/:app/*path", func(ctx *gin.Context) {
		appName := ctx.Param("app")
		path := ctx.Param("path")

		reqBody, reqBodyErr := ioutil.ReadAll(ctx.Request.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		loadRequestsFromFile(appName)
		storedRes := findResponseByRequestParams(ctx.Request.Method, path, reqBody)
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

			storeResponseAndRequest(appName, ctx.Request, path, reqBody, cRes, cBody)
		}
	})

	s := &http.Server{
		Addr:           *addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
