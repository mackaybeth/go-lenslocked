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
	r := chi.NewRouter()

	// layout-page must be first because the page template wraps everything in home.gohtml
	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(faqTpl))

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

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

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

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)
	// Annoying to create links and forms that peform DELETE without JS, so we're using POST
	r.Post("/signout", usersC.ProcessSignOut)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound)+": "+r.URL.Path, http.StatusNotFound)
	})

	// Create an instance of the user middleware
	usrMw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	var csrfKey = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // 32-byte key
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false)) // TODO Fix this before deploy

	fmt.Println("Starting the server on :3000...")

	// Note on the nested MW calls around r:  usrMw returns a new request, then the result of that is passed into csrfMw which returns another new request
	// Ordering is important.  csrfMw is first, then usrMw is next, then that wrapped result is sent into ListenAndServe
	http.ListenAndServe("localhost:3000", csrfMw(usrMw.SetUser(r)))

}
