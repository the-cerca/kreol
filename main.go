package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"kreol/controllers"
	"kreol/models"
	"kreol/statics"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "djscqdnsqdlfknlqsfkfnqsldnfqsief ezifazi")
	if err != nil {
		fmt.Printf("Error connecting database %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		fmt.Printf("Ping database %v", err)
	}
	mw := controllers.Middleware{
		Sm: &models.SessionManager{DB: db},
	} 
	uc := controllers.UserController{
		Um: &models.UserManager{DB: db},
		Sm: &models.SessionManager{DB: db},
	}
	lc := controllers.LanguageController{
		Lm: &models.LanguageManager{DB: db},
	}
	r := chi.NewRouter()

	// server static file
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(statics.StaticsFs))))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFS(statics.StaticsFs, "views/404.gohtml", "views/tools.gohtml", "views/navbar.gohtml"))
		t.Execute(w, nil)
	})
	//public route
	r.Group(func(r chi.Router) {
		r.Get("/", uc.HandleGetHome)
		r.Get("/signin", uc.HandleGetSignIn)
		r.Post("/signin", uc.HandlePostSignIn)
		r.Get("/login", uc.HandleGetLogin)
		r.Post("/login", uc.HandlePostLogin)
	})
	// private route
	r.Route("/lang", func(r chi.Router) {
		r.With(mw.GetCurrentUser).Get("/selection", lc.HandleGetLanguages)
		r.With(mw.GetCurrentUser).Post("/selection/{id}", lc.HandleSubscribeLanguage)
		r.With(mw.GetCurrentUser).Get("/{lang}", lc.HandleGetCourseLanguage)
		r.With(mw.GetCurrentUser).Post("/play/{id}", lc.PlayReponse)
	})
	log.Println("http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
