package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func getUrlParts(regex *regexp.Regexp, url *url.URL) (string, string) {
	appName, pUrl := getUrl(regex, url)
	pUrl = "/" + pUrl

	return appName, pUrl
}

func getUrl(regex *regexp.Regexp, rUrl *url.URL) (string, string) {
	matches := regex.FindStringSubmatch(rUrl.String())
	return matches[1], matches[2]
}

func createClientRequest(method string, url string, body []byte) *http.Request {
	r, rErr := http.NewRequest(method, url, bytes.NewReader(body))
	if rErr != nil {
		log.Fatal(rErr)
	}

	return r
}

func doClientRequest(c *http.Client, r *http.Request) *http.Response {
	res, resErr := c.Do(r)
	if resErr != nil {
		log.Fatal(resErr)
	}

	return res
}

func copyHeaders(headers http.Header, ctx *gin.Context) {
	for hKey, hValues := range headers {
		for _, hValue := range hValues {
			ctx.Header(hKey, hValue)
		}
	}
}

func copyBodyToClient(w http.ResponseWriter, cBody []byte) {
	_, wErr := w.Write(cBody)
	if wErr != nil {
		log.Fatal(wErr)
	}
}

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
			copyHeaders(storedRes.Header, ctx)
			ctx.Writer.WriteHeader(storedRes.StatusCode)
			copyBodyToClient(ctx.Writer, storedRes.Body)
		} else {
			fmt.Println("Proxying response to real backend")
			cReq := createClientRequest(ctx.Request.Method, "http://host.docker.internal:4000"+path, reqBody)
			defer cReq.Body.Close()

			cRes := doClientRequest(client, cReq)

			cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
			if cBodyErr != nil {
				log.Fatal(cBodyErr)
			}

			fmt.Printf("Client response body: %v\n", string(cBody))

			copyHeaders(cReq.Header, ctx)
			ctx.Writer.WriteHeader(cRes.StatusCode)
			copyBodyToClient(ctx.Writer, cBody)

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
