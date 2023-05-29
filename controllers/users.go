package controllers

import (
	"fmt"
	"net/http"

	"github.com/mackaybeth/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, data)
}

// Parsing the values from the form using helper methosds in ther request
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// These are getting the "name" in the html (not the id or type, even though they are all named the same)
	// FormValue automatically parses the form, so no need to call the functions to do that
	// FormValue does not return errors though, so if you need to parse the error then you need to use the other method
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User crated: %+v", user)
}
