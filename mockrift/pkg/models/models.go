package models

import (
	"bytes"
	"fmt"
	"net/http"
)

type App struct {
	Name     string
	Alias    string
	Requests []*StoredRequest
}

type StoredResponse struct {
	Active     bool        `file:"active"`
	StatusCode int         `file:"status_code"`
	Header     http.Header `file:"header"`
	Body       []byte      `file:"body"`
}

type StoredRequest struct {
	Method    string            `file:"method"`
	URL       string            `file:"url"`
	Header    http.Header       `file:"header"`
	Body      []byte            `file:"body"`
	Responses []*StoredResponse `file:"responses"`
}

func (a *App) GetName() string {
	if a.Alias != "" {
		return a.Alias
	}

	return a.Name
}

func (a *App) FindRequest(m string, url string, body []byte) *StoredRequest {
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

func (a *App) FindResponseByRequestParams(m string, url string, body []byte) *StoredResponse {
	r := a.FindRequest(m, url, body)
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

func (a *App) AddResponseAndRequest(req *http.Request, url string, reqBody []byte, res *http.Response, resBody []byte) {
	sRes := StoredResponse{
		Active:     false,
		Header:     res.Header,
		StatusCode: res.StatusCode,
		Body:       resBody,
	}

	r := a.FindRequest(req.Method, url, reqBody)
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
}
