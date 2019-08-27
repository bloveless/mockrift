package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type StoredResponse struct {
	Body []byte
	Header http.Header
	StatusCode int
}

type StoredRequest struct {
	Method string
	URL string
	Body []byte
	Header http.Header
	Response *StoredResponse
}

var requests []*StoredRequest

func findResponseByRequestParams(m string, url string, body []byte) *StoredResponse {
	for _, r := range requests {
		if r.Method == m && r.URL == url {
			fmt.Printf("Matched method (%s) and url (%s)\n", m, url)
			if bytes.Compare(r.Body, body) == 0 {
				fmt.Println("Matched body")
				return r.Response
			}
		}
	}

	return nil
}

func storeResponseAndRequest(req *http.Request, reqBody []byte, res *http.Response, resBody []byte) {
	sRes := StoredResponse{
		Body:       resBody,
		Header:     res.Header,
		StatusCode: res.StatusCode,
	}

	sReq := StoredRequest{
		Method:   req.Method,
		URL:      req.URL.String(),
		Body:     reqBody,
		Response: &sRes,
	}

	requests = append(requests, &sReq)
}
