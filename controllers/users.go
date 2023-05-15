package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// we need a view to render
	u.Templates.New.Execute(w, nil)
}

// Parsing the values from the form using helper methosds in ther request
func (u Users) Create(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// These are getting the "name" in the html (not the id or type, even though they are all named the same)
	fmt.Fprint(w, "Email: ", r.PostForm.Get("email"))
	fmt.Fprint(w, "Password: ", r.PostForm.Get("password"))
}
