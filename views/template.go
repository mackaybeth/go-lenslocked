package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
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
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			// Name of the function : type returnval
			"csrfField": func() template.HTML {
				// This is a placeholder, we want to use csrf.TemplateField but
				// requests are not available here.  This allows us to put in the placeholder
				// and later update the template to replace the csrfField function
				return `<!-- TODO: Implement the csrfField func -->`
			},
		},
	)

	// Need to add the 3 dots after the input patterns (even though both take in
	// variadic string) to tell template.ParseFS to treat this slice as a variadic string
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parseFS template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil

}

// func Parse(filepath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filepath)
// 	if err != nil {
// 		// log.Printf("parsing the template: %v", err)
// 		// http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}
// 	return Template{
// 		htmlTpl: tpl,
// 	}, nil
// }

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {

	// t.htmlTpl is a pointer, so we want to clone before using to avoid race conditions
	// where multiple requests come in at once and all update the same pointer
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			// Name of the function : type returnval
			"csrfField": func() template.HTML {
				// This is ovewriting what we parsed originally in ParseFS
				// (because now we have an http.Request)
				return csrf.TemplateField(r)
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Pass in the http.ResponseWriter as the place to write the template
	err = tpl.Execute(w, data)

	// update the functions specific to the request
	//t.htmlTpl.Funcs()
	if err != nil {
		log.Printf("executing the template: %v", err)
		// This doesn't actually work to set an error, because when the template executes it starts
		// rendering things (and sets the response to 200 which can't be changed).  We see valid data
		// and then an error on the resulting page, not a dedicated error page.  This is expected.
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
