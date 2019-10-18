package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/graphql-go/handler"
	"html/template"
	"log"
	"mockrift/pkg/models/file"
	"net/http"
	"os"
	"time"
)

type server struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	client        *http.Client
	recordOnly    bool
	apps          *file.AppModel
	templateCache map[string]*template.Template
	templateData  *templateData
}

func main() {
	addr := flag.String("addr", ":3499", "The address to run the mockrift server")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	tc, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	td := &templateData{}
	addDefaultData(td)

	s := &server{
		infoLog:       infoLog,
		errorLog:      errorLog,
		client:        client,
		recordOnly:    false,
		apps:          &file.AppModel{},
		templateCache: tc,
		templateData:  td,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Default().Handler)

	fs := http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static")))

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	schema := s.getSchema()

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	r.Handle("/admin/graphql", h)
	r.Mount("/admin", s.adminRouter())
	r.Mount("/m", s.mockRouter())

	httpServer := &http.Server{
		Addr:           *addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Starting server at %s\n", *addr)
	log.Fatal(httpServer.ListenAndServe())
}
