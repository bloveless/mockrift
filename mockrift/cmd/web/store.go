package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type App struct {
	Name     string
	Requests []*StoredRequest
}

type StoredResponse struct {
	Active     bool        `json:"active"`
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

type StoredRequest struct {
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Header    http.Header       `json:"header"`
	Body      []byte            `json:"body"`
	Responses []*StoredResponse `json:"responses"`
}

func getAllApps() []*App {
	appFiles, gErr := filepath.Glob("./requests/*.json")
	if gErr != nil {
		log.Fatal("Error loading request files: " + gErr.Error())
	}

	var apps []*App

	for _, appFile := range appFiles {
		appBytes, rErr := ioutil.ReadFile(appFile)
		if rErr != nil {
			log.Fatal("Unable to read app file: " + rErr.Error())
		}

		var app App
		uErr := json.Unmarshal(appBytes, &app)
		if uErr != nil {
			log.Fatal("Unable to unmarshal app json file: " + uErr.Error())
		}

		apps = append(apps, &app)
	}

	return apps
}

// getAppFromFile will load any stored requests for an app prefix.
func getAppFromFile(appName string) *App {
	var a App

	fmt.Println("Loading app from /home/appuser/app/requests/" + appName + ".json")
	jsonFile, err := os.Open("/home/appuser/app/requests/" + appName + ".json")
	if err != nil {
		// If the file doesn't exist then that is fine. We'll just save the file upon the first response.
		return nil
	}
	defer jsonFile.Close()

	jsonBytes, jsonBytesErr := ioutil.ReadAll(jsonFile)
	if jsonBytesErr != nil {
		log.Fatal(fmt.Printf("Unable to read JSON file (%s): %s\n", appName, jsonBytesErr.Error()))
	}

	unmarshalErr := json.Unmarshal(jsonBytes, &a)
	if unmarshalErr != nil {
		log.Fatal("Unable to unmarshal json file: " + unmarshalErr.Error())
	}

	return &a
}

func (a *App) findRequest(m string, url string, body []byte) *StoredRequest {
	for _, r := range a.Requests {
		if r.Method == m && r.URL == url {
			fmt.Printf("Matched method (%s) and url (%s)\n", m, url)
			if bytes.Compare(r.Body, body) == 0 {
				fmt.Println("Matched body")
				return r
			}
		}
	}

	return nil
}

func (a *App) findResponseByRequestParams(m string, url string, body []byte) *StoredResponse {
	r := a.findRequest(m, url, body)
	if r == nil {
		return nil
	}

	for _, response := range r.Responses {
		if response.Active == true {
			return response
		}
	}

	return r.Responses[0]
}

func (a *App) storeResponseAndRequest(req *http.Request, url string, reqBody []byte, res *http.Response, resBody []byte) {
	sRes := StoredResponse{
		Active:     false,
		Header:     res.Header,
		StatusCode: res.StatusCode,
		Body:       resBody,
	}

	r := a.findRequest(req.Method, url, reqBody)
	if r != nil {
		if len(r.Responses) == 0 {
			sRes.Active = true
		}

		r.Responses = append(r.Responses, &sRes)
	} else {
		sRes.Active = true

		sReq := StoredRequest{
			Method: req.Method,
			URL:    url,
			Header: req.Header,
			Body:   reqBody,
			Responses: []*StoredResponse{
				&sRes,
			},
		}

		a.Requests = append(a.Requests, &sReq)
	}

	requestsJson, marshalErr := json.MarshalIndent(a.Requests, "", "  ")
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	f, oErr := os.OpenFile("./requests/"+a.Name+".json", os.O_WRONLY|os.O_CREATE, 0644)
	if oErr != nil {
		log.Fatal("Unable to open file for writing: " + oErr.Error())
	}

	_, wErr := f.Write(requestsJson)
	if wErr != nil {
		log.Fatal("Unable to write json to file: " + wErr.Error())
	}
}
