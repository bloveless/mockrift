package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-http-utils/logger"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

func copyHeaders(headers http.Header, w http.ResponseWriter) {
	for hKey, hValues := range headers {
		for _, hValue := range hValues {
			w.Header().Add(hKey, hValue)
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

	urlRegex := regexp.MustCompile("/(.+?)/(.*)")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	recordOnly := false

	mux := http.NewServeMux()

	tc, newTcErr := newTemplateCache("./ui/html")
	if newTcErr != nil {
		log.Fatal("Unable to create template cache: " + newTcErr.Error())
	}

	p := &Pages{
		templateCache: tc,
	}

	mux.HandleFunc("/admin", p.handleHome)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		appName, pUrl := getUrlParts(urlRegex, req.URL)

		fmt.Println("Handling rest")

		reqBody, reqBodyErr := ioutil.ReadAll(req.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		loadRequestsFromFile(appName)
		storedRes := findResponseByRequestParams(req.Method, pUrl, reqBody)
		if !recordOnly && storedRes != nil {
			fmt.Println("Sending response from memory")
			copyHeaders(storedRes.Header, w)
			w.WriteHeader(storedRes.StatusCode)
			copyBodyToClient(w, storedRes.Body)
		} else {
			fmt.Println("Proxying response to real backend")
			cReq := createClientRequest(req.Method, "http://host.docker.internal:4000"+pUrl, reqBody)
			defer cReq.Body.Close()

			cRes := doClientRequest(client, cReq)

			cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
			if cBodyErr != nil {
				log.Fatal(cBodyErr)
			}

			fmt.Printf("Client response body: %v\n", string(cBody))

			copyHeaders(cReq.Header, w)
			w.WriteHeader(cRes.StatusCode)
			copyBodyToClient(w, cBody)

			storeResponseAndRequest(appName, req, pUrl, reqBody, cRes, cBody)
		}
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	corsHandler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(*addr, logger.Handler(corsHandler, os.Stdout, logger.DevLoggerType)))
}
