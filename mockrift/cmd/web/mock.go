package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"mockrift/pkg/models"
	"net/http"
)

func (s *server) mockRouter() http.Handler {
	r := chi.NewRouter()
	r.Handle("/{app:[s-z-]+}/{path}", s.handleMock())
	return r
}

func (s *server) handleMock() http.HandlerFunc {
	var app *models.App

	return func(w http.ResponseWriter, r *http.Request) {
		appName := chi.URLParam(r, "app")
		path := "/" + chi.URLParam(r, "path")

		reqBody, reqBodyErr := ioutil.ReadAll(r.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		if app == nil || app.Name != appName {
			app = s.apps.Get(appName)
		}

		storedRes := app.FindResponseByRequestParams(r.Method, path, reqBody)
		if !s.recordOnly && storedRes != nil {
			fmt.Println("Sending response from memory")
			copyHeadersFromStoredResponse(storedRes.Header, w)
			w.WriteHeader(storedRes.StatusCode)
			copyBodyToClient(w, storedRes.Body)
		} else {
			fmt.Println("Proxying response to real backend")
			cReq := createClientRequest(r.Method, "http://host.docker.internal:4000"+path, reqBody)
			defer cReq.Body.Close()

			cRes := doClientRequest(s.client, cReq)

			cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
			if cBodyErr != nil {
				log.Fatal(cBodyErr)
			}

			fmt.Printf("Client response body: %v\n", string(cBody))

			copyHeadersFromRequest(cReq.Header, w)
			w.WriteHeader(cRes.StatusCode)
			copyBodyToClient(w, cBody)

			app.AddResponseAndRequest(r, path, reqBody, cRes, cBody)
			s.apps.Save(app)
		}
	}
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

func copyHeadersFromStoredResponse(headers []*models.StoredHeader, w http.ResponseWriter) {
	for _, header := range headers {
		for _, hValue := range header.Value {
			w.Header().Add(header.Name, hValue)
		}
	}
}

func copyHeadersFromRequest(headers http.Header, w http.ResponseWriter) {
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
