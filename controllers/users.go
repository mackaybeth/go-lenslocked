package controllers

import (
	"fmt"
	"net/http"

	"github.com/mackaybeth/lenslocked/context"
	"github.com/mackaybeth/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
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
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: long term, we should show a warning about not being able to sign the user in
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)

	// Redirect to the current user page
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Somethng when wrong.", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// This is from the context package that we wrote
	user := context.User(ctx)
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user %s\n", user.Email)

	// token, err := readCookie(r, CookieSession)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	// user, err := u.SessionService.User(token)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }

	// fmt.Fprintf(w, "Current user %s\n", user.Email)
	// fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add logic for the SetUser middleware, then eventually call next.ServeHTTP(w, r)

		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// At this point we have a user
		ctx := r.Context()
		// This is the context package we wrote
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
