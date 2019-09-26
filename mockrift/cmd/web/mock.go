package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"mockrift/pkg/helper"
	"mockrift/pkg/models"
	"net/http"
)

func (a *application) mockRouter() http.Handler {
	r := chi.NewRouter()
	r.Handle("/{app:[a-z-]+}/{path}", a.handleMock())
	return r
}

func (a *application) handleMock() http.HandlerFunc {
	var app *models.App

	return func(w http.ResponseWriter, r *http.Request) {
		appName := chi.URLParam(r, "app")
		path := chi.URLParam(r, "path")

		reqBody, reqBodyErr := ioutil.ReadAll(r.Body)
		if reqBodyErr != nil {
			log.Fatal("Unable to read request body")
		}

		if app == nil || app.Name != appName {
			app = a.apps.Get(appName)
		}

		storedRes := app.FindResponseByRequestParams(r.Method, path, reqBody)
		if !a.recordOnly && storedRes != nil {
			fmt.Println("Sending response from memory")
			helper.CopyHeaders(storedRes.Header, w)
			w.WriteHeader(storedRes.StatusCode)
			helper.CopyBodyToClient(w, storedRes.Body)
		} else {
			fmt.Println("Proxying response to real backend")
			cReq := helper.CreateClientRequest(r.Method, "http://host.docker.internal:4000"+path, reqBody)
			defer cReq.Body.Close()

			cRes := helper.DoClientRequest(a.client, cReq)

			cBody, cBodyErr := ioutil.ReadAll(cRes.Body)
			if cBodyErr != nil {
				log.Fatal(cBodyErr)
			}

			fmt.Printf("Client response body: %v\n", string(cBody))

			helper.CopyHeaders(cReq.Header, w)
			w.WriteHeader(cRes.StatusCode)
			helper.CopyBodyToClient(w, cBody)

			app.AddResponseAndRequest(r, path, reqBody, cRes, cBody)
		}
	}
}
