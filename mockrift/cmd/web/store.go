package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

var currentApp string
var requests []*StoredRequest

func loadRequestsFromFile(app string) {
	if currentApp == app {
		return
	}

	fmt.Println("Loading stored requests from /home/appuser/app/requests/" + app + ".json")
	jsonFile, err := os.Open("/home/appuser/app/requests/" + app + ".json")
	if err != nil {
		// If the file doesn't exist then that is fine. We'll just save the file upon the first response.
		return
	}
	defer jsonFile.Close()

	jsonBytes, jsonBytesErr := ioutil.ReadAll(jsonFile)
	if jsonBytesErr != nil {
		log.Fatal(fmt.Printf("Unable to read JSON file (%s): %s\n", app, jsonBytesErr.Error()))
	}

	unmarshalErr := json.Unmarshal(jsonBytes, &requests)
	if unmarshalErr != nil {
		log.Fatal("Unable to unmarshal json file: " + unmarshalErr.Error())
	}

	currentApp = app
}

func findRequest(m string, url string, body []byte) *StoredRequest {
	for _, r := range requests {
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

func findResponseByRequestParams(m string, url string, body []byte) *StoredResponse {
	r := findRequest(m, url, body)
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

func storeResponseAndRequest(app string, req *http.Request, url string, reqBody []byte, res *http.Response, resBody []byte) {
	sRes := StoredResponse{
		Active:     false,
		Header:     res.Header,
		StatusCode: res.StatusCode,
		Body:       resBody,
	}

	r := findRequest(req.Method, url, reqBody)
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

		requests = append(requests, &sReq)
	}

	requestsJson, marshalErr := json.MarshalIndent(requests, "", "  ")
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	writeErr := ioutil.WriteFile("/home/appuser/app/requests/"+app+".json", requestsJson, 0644)
	if writeErr != nil {
		log.Fatal("Unable to write json to file: " + writeErr.Error())
	}
}
