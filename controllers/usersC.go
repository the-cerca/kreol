package controllers

import (
	"html/template"
	"kreol/models"
	"kreol/statics"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

type UserController struct {
	Um *models.UserManager
	Sm *models.SessionManager
}
type CustomErrors struct {
	Errors map[string]string
}

func (uc *UserController) HandleGetHome(w http.ResponseWriter, _ *http.Request) {
	v := AddViews("views/home.gohtml")
	t := template.Must(template.ParseFS(statics.StaticsFs, v...))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
func (uc *UserController) HandleGetSignIn(w http.ResponseWriter, _ *http.Request) {
	v := AddViews("views/sign.gohtml")
	t := template.Must(template.ParseFS(statics.StaticsFs, v...))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
func (uc *UserController) HandlePostSignIn(w http.ResponseWriter, r *http.Request) {
	v := AddViews("views/sign.gohtml")
	users := models.UserCreate{
		Username:       strings.TrimSpace(r.FormValue("username")),
		Email:          strings.TrimSpace(strings.ToLower(r.FormValue("email"))),
		Password:       r.FormValue("password"),
		RepeatPassword: r.FormValue("repeat-password"),
	}
	errors := make(map[string]string)
	if users.RepeatPassword != users.Password {
		errors["password"] = "Mots de passe non identiques."
	}
	if _, err := mail.ParseAddress(users.Email); err != nil {
		errors["email"] = "Adresse email incorrect."
	}
	if b := uc.Um.SearchUserByEmail(r.Context(), users.Email); b {
		errors["email"] = "Adresse déjà utilisé."
	}
	if len(errors) > 0 {
		data := CustomErrors{
			Errors: errors,
		}
		t := template.Must(template.ParseFS(statics.StaticsFs, v...))
		_ = t.Execute(w, data)
		return
	}

	user, err := uc.Um.Create(r.Context(), users.Username, users.Email, users.Password)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		log.Println("create error :", err)
		return
	}

	if user != nil {
		s, err := uc.Sm.CreateSession(r.Context(), user.ID)
		if err != nil {
			log.Printf("error %s", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		SetCookie(w, "session_id", s.Token)
		http.Redirect(w, r, "/lang/selection", http.StatusMovedPermanently)
	}
}
func (uc *UserController) HandleGetLogin(w http.ResponseWriter, _ *http.Request) {
	v := AddViews("views/login.gohtml")
	t := template.Must(template.ParseFS(statics.StaticsFs, v...))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (uc *UserController) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	v := AddViews("views/login.gohtml")
	e := strings.ToLower(r.FormValue("email"))
	p := r.FormValue("password")
	id, err := uc.Um.Login(r.Context(), e, p)
	if err != nil {
		data := CustomErrors{
			Errors: map[string]string{"password": models.ErrorPasswordEmailInvalid.Error()},
		}
		t := template.Must(template.ParseFS(statics.StaticsFs, v...))
		err := t.Execute(w, data)
		if err != nil {
			return
		}
	}
	s, err := uc.Sm.CreateSession(r.Context(), id)
	if err != nil {
		log.Printf("Create session : %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	SetCookie(w, "session_id", s.Token)
	http.Redirect(w, r, "/selection", http.StatusMovedPermanently)
}
func (uc *UserController) HandleGetHomeUser(w http.ResponseWriter, _ *http.Request) {
	v := AddViews("views/play.gohtml")
	t := template.Must(template.ParseFS(statics.StaticsFs, v...))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

func AddViews(views ...string) []string {
	base := []string{"views/tools.gohtml"}
	return append(views, base...)
}
