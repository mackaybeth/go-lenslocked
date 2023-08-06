package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/mackaybeth/lenslocked/context"
	"github.com/mackaybeth/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

// We will use this to determine if an error provides the Public method.
type public interface {
	Public() string
}

// Helper function used for templates, to wrap the panic
// Major benefit is to reduce copypasta in main
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0]))
	tpl = tpl.Funcs(
		template.FuncMap{
			// Name of the function : type returnval
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented, add code to Execute to implement")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented, add code to Execute to implement")
			},
			"errors": func() []string {
				return nil
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {

	// t.htmlTpl is a pointer, so we want to clone before using to avoid race conditions
	// where multiple requests come in at once and all update the same pointer
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	// Call the errMessages func before the closures.
	errMsgs := errMessages(errs...)

	tpl = tpl.Funcs(
		template.FuncMap{
			// Name of the function : type returnval
			"csrfField": func() template.HTML {
				// This is ovewriting what we parsed originally in ParseFS
				// (because now we have an http.Request)
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer

	// Pass in the bytes.Buffer as the place to write the template, in case
	// there is an error
	err = tpl.Execute(&buf, data)

	// update the functions specific to the request
	if err != nil {
		log.Printf("executing the template: %v", err)
		// This doesn't actually work to set an error, because when the template executes it starts
		// rendering things (and sets the response to 200 which can't be changed).  We see valid data
		// and then an error on the resulting page, not a dedicated error page.  This is expected.
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	// We got this far, so no error.  Now we can copy from the buffer to the
	// http.ResponseWriter
	io.Copy(w, &buf)
}
