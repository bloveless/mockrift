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

	fmt.Println(a.templateData)

	execErr := tmpl.Execute(w, a.templateData)
	if execErr != nil {
		a.serverError(w, fmt.Errorf("unable to execute template: %s", execErr.Error()))
	}
}
