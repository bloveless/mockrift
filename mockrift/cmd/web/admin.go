package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func (a *application) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/*", a.handleReact)
	return r
}

func (a *application) handleReact(w http.ResponseWriter, r *http.Request) {
	tmpl := a.templateCache["react.tmpl"]

	execErr := tmpl.Execute(w, nil)
	if execErr != nil {
		a.serverError(w, fmt.Errorf("unable to execute template: %s", execErr.Error()))
	}
}
