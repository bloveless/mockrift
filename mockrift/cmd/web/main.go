package main

import (
	"flag"
	"fmt"
	"github.com/go-http-utils/logger"
	"github.com/rs/cors"
	"io"
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

func createClientRequest(method string, url string, body io.Reader) *http.Request {
	r, rErr := http.NewRequest(method, url, body)
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

func copyHeadersToClient(w http.ResponseWriter, cReq *http.Request, cRes *http.Response) {
	for hKey, hValues := range w.Header() {
		for _, hValue := range hValues {
			cReq.Header.Add(hKey, hValue)
		}
	}
	w.WriteHeader(cRes.StatusCode)
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
		cReq := createClientRequest(req.Method, "http://host.docker.internal:4000"+pUrl, req.Body)
		defer cReq.Body.Close()

		cRes := doClientRequest(client, cReq)

		cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
		if cBodyErr != nil {
			log.Fatal(cBodyErr)
		}

		fmt.Printf("Client Response body: %v\n", string(cBody))

		copyHeadersToClient(w, cReq, cRes)
		copyBodyToClient(w, cBody)
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	corsHandler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(*addr, logger.Handler(corsHandler, os.Stdout, logger.DevLoggerType)))
}
