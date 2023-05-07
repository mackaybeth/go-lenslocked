package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mackaybeth/lenslocked/controllers"
	"github.com/mackaybeth/lenslocked/templates"
	"github.com/mackaybeth/lenslocked/views"
)

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found</h1><p>Path not supported: "+r.URL.Path)
}

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout-parts.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "layout-parts.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "layout-parts.gohtml"))
	r.Get("/faq", controllers.FAQ(faqTpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	// http.HandlerFunc is a type conversion,  NOT a funciton call
	http.ListenAndServe("localhost:3000", r)
}
