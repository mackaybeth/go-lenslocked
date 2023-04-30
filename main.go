package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/mackaybeth/lenslocked/views"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	// tpl, err := template.ParseFiles(filepath)
	// if err != nil {
	// 	log.Printf("parsing the template: %v", err)
	// 	http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
	// 	return
	// }
	// // 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 	// Pass in the http.ResponseWriter as the place to write the template
	// 	err = tpl.Execute(w, nil)
	// 	if err != nil {
	// 		log.Printf("executing the template: %v", err)
	// 		// This doesn't actually work to set an error, because when the template executes it starts
	// 		// rendering things (and sets the response to 200 which can't be changed).  We see valid data
	// 		// and then an error on the resulting page, not a dedicated error page.  This is expected.
	// 		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
	// 		return
	// 	}

	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing the template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml") // makes the path os-agnostic
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found</h1><p>Path not supported: "+r.URL.Path)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	// http.HandlerFunc is a type conversion,  NOT a funciton call
	http.ListenAndServe("localhost:3000", r)
}
