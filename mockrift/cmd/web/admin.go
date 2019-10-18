package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func (s *server) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/*", s.handleReact)
	return r
}

func (s *server) handleReact(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templateCache["react.tmpl"]

	fmt.Println(s.templateData)

	execErr := tmpl.Execute(w, s.templateData)
	if execErr != nil {
		s.serverError(w, fmt.Errorf("unable to execute template: %s", execErr.Error()))
	}
}
