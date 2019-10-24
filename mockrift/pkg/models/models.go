package models

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type App struct {
	Slug     string           `json:"slug"`
	Name     string           `json:"name"`
	Requests []*StoredRequest `json:"requests"`
}

type StoredHeader struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type StoredResponse struct {
	ID         string          `json:"id"`
	Active     bool            `json:"active"`
	StatusCode int             `json:"status_code"`
	Header     []*StoredHeader `json:"header"`
	Body       []byte          `json:"body"`
}

type StoredRequest struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Header    []*StoredHeader   `json:"header"`
	Body      []byte            `json:"body"`
	Responses []*StoredResponse `json:"responses"`
}

func (a *App) GetName() string {
	if a.Name != "" {
		return a.Name
	}

	return a.Slug
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
		ID:       uuid.New().String(),
		Active:     false,
		StatusCode: res.StatusCode,
		Body:       resBody,
	}

	var resHeaders []*StoredHeader
	for name, value := range res.Header {
		resHeaders = append(resHeaders, &StoredHeader{
			Name:  name,
			Value: value,
		})
	}

	sRes.Header = resHeaders

	r := a.FindRequest(req.Method, url, reqBody)
	if r != nil {
		if len(r.Responses) == 0 {
			sRes.Active = true
		}

		r.Responses = append(r.Responses, &sRes)
	} else {
		sRes.Active = true

		sReq := StoredRequest{
			ID:   uuid.New().String(),
			Method: req.Method,
			URL:    url,
			Body:   reqBody,
			Responses: []*StoredResponse{
				&sRes,
			},
		}

		var reqHeaders []*StoredHeader
		for name, value := range req.Header {
			reqHeaders = append(reqHeaders, &StoredHeader{
				Name:  name,
				Value: value,
			})
		}

		sReq.Header = reqHeaders

		a.Requests = append(a.Requests, &sReq)
	}
}
