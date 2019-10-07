package helper

// TODO: This should not stay here. This is here just to experiment with packages.

import (
	"bytes"
	"log"
	"mockrift/pkg/models"
	"net/http"
	"net/url"
	"regexp"
)

func GetUrl(regex *regexp.Regexp, rUrl *url.URL) (string, string) {
	matches := regex.FindStringSubmatch(rUrl.String())
	return matches[1], matches[2]
}

func CreateClientRequest(method string, url string, body []byte) *http.Request {
	r, rErr := http.NewRequest(method, url, bytes.NewReader(body))
	if rErr != nil {
		log.Fatal(rErr)
	}

	return r
}

func DoClientRequest(c *http.Client, r *http.Request) *http.Response {
	res, resErr := c.Do(r)
	if resErr != nil {
		log.Fatal(resErr)
	}

	return res
}

func CopyHeadersFromStoredResponse(headers []*models.StoredHeader, w http.ResponseWriter) {
	for _, header := range headers {
		for _, hValue := range header.Value {
			w.Header().Add(header.Name, hValue)
		}
	}
}

func CopyHeadersFromRequest(headers http.Header, w http.ResponseWriter) {
	for hKey, hValues := range headers {
		for _, hValue := range hValues {
			w.Header().Add(hKey, hValue)
		}
	}
}

func CopyBodyToClient(w http.ResponseWriter, cBody []byte) {
	_, wErr := w.Write(cBody)
	if wErr != nil {
		log.Fatal(wErr)
	}
}
