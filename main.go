package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mackaybeth/lenslocked/controllers"
	"github.com/mackaybeth/lenslocked/models"
	"github.com/mackaybeth/lenslocked/templates"
	"github.com/mackaybeth/lenslocked/views"
)

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found</h1><p>Path not supported: "+r.URL.Path)
}

func main() {
	r := chi.NewRouter()

	// layout-page must be first because the page template wraps everything in home.gohtml
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(faqTpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService, // takes a pointer
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.New)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	// http.HandlerFunc is a type conversion,  NOT a funciton call
	http.ListenAndServe("localhost:3000", r)
}
