package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (s *server) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *server) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (s *server) notFound(w http.ResponseWriter) {
	s.clientError(w, http.StatusNotFound)
}
