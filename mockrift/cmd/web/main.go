package main

import (
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", ":3499", "The address to run the mockrift server")
	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	recordOnly := false

	router := gin.Default()
	router.Use(cors.Default())
	router.LoadHTMLGlob("./ui/html/*")

	router.StaticFS("/static", http.Dir("./ui/static"))

	addAdminRoutes(router)
	addMockingRoutes(router, client, recordOnly)

	s := &http.Server{
		Addr:           *addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
