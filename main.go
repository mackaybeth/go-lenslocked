package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/mackaybeth/lenslocked/controllers"
	"github.com/mackaybeth/lenslocked/migrations"
	"github.com/mackaybeth/lenslocked/models"
	"github.com/mackaybeth/lenslocked/templates"
	"github.com/mackaybeth/lenslocked/views"
)

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found</h1><p>Path not supported: "+r.URL.Path)
}

func main() {

	// SETUP THE DATABASE
	cfg := models.DefaultPostgresConfig()

	// Print out the config for the DB so we can use DB migrations
	fmt.Println(cfg.String())

	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// SETUP SERVICES

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// SETUP MIDDLEWARE
	usrMw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	var csrfKey = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // 32-byte key
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false)) // TODO Fix this before deploy

	// SETUP CONTROLLERS
	usersC := controllers.Users{
		UserService:    &userService, // takes a pointer
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	// SETUP ROUTER AND ROUTES

	r := chi.NewRouter()

	//  Make the router use the middlewre
	r.Use(csrfMw)
	r.Use(usrMw.SetUser)

	// layout-page must be first because the page template wraps everything in home.gohtml
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(faqTpl))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	// Annoying to create links and forms that peform DELETE without JS, so we're using POST
	r.Post("/signout", usersC.ProcessSignOut)

	// Can use subroute "Route" here because we know that this prefix means that user needs to be logged in
	r.Route("/users/me", func(r chi.Router) {
		// This MW is used by pages under the /users/me prefix
		r.Use(usrMw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "howdy")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})

	// START THE SERVER
	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe("localhost:3000", r)

}
