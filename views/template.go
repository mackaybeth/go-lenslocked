package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

// Helper function used for templates, to wrap the panic
// Major benefit is to reduce copypasta in main
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {

	// Need to add the 3 dots after the input patterns (even though both take in
	// variadic string) to tell template.ParseFS to treat this slice as a variadic string
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parseFS template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil

}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		// log.Printf("parsing the template: %v", err)
		// http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Pass in the http.ResponseWriter as the place to write the template
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing the template: %v", err)
		// This doesn't actually work to set an error, because when the template executes it starts
		// rendering things (and sets the response to 200 which can't be changed).  We see valid data
		// and then an error on the resulting page, not a dedicated error page.  This is expected.
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
