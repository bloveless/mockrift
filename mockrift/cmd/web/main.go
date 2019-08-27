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

func getUrl(regex *regexp.Regexp, rUrl *url.URL) string {
	return string(regex.ReplaceAll([]byte(rUrl.String()), []byte("/$2")))
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

	urlRegex, urlRegexErr := regexp.Compile("/(.+?)/(.*)")
	if urlRegexErr != nil {
		log.Fatal(urlRegexErr)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		pUrl := getUrl(urlRegex, req.URL)
		reqBody, reqBodyErr := ioutil.ReadAll(req.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		storedRes := findResponseByRequestParams(req.Method, req.URL.String(), reqBody)
		if storedRes != nil {
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

			storeResponseAndRequest(req, reqBody, cRes, cBody)
		}
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	corsHandler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(*addr, logger.Handler(corsHandler, os.Stdout, logger.DevLoggerType)))
}
