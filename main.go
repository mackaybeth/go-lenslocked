package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/mackaybeth/lenslocked/controllers"
	"github.com/mackaybeth/lenslocked/views"
)

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found</h1><p>Path not supported: "+r.URL.Path)
}

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	r.Get("/", controllers.StaticHandler(tpl))

	contactTpl := views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	faqTpl := views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))
	r.Get("/faq", controllers.StaticHandler(faqTpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	// http.HandlerFunc is a type conversion,  NOT a funciton call
	http.ListenAndServe("localhost:3000", r)
}
