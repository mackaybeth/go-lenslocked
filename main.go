package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplPath := filepath.Join("templates", "home.gohtml") // makes the path os-agnostic
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("parsing the template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	// Pass in the http.ResponseWriter as the place to write the template
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing the template: %v", err)
		// This doesn't actually work to set an error, because when the template executes it starts
		// rendering things (and sets the response to 200 which can't be changed).  We see valid data
		// and then an error on the resulting page, not a dedicated error page.  This is expected.
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:email@email.com\">email@email.com</a>.")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Frequently Asked Questions</h1>
	<h2>Q: Is there a free version?</h2>
	<p><b>A:</b> Yes! We offer a free trial for 30 days</p>
	<h2>Q: What are your support hours?</h2>
	<p><b>A:</b> 24/7 email support, slower on weekends</p>
	<h2>Q: How do Io contact support?</h2>
	<p><b>A:</b> email <a href="mailto:support@email.com">support@email.com</a></p>
	`)
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
