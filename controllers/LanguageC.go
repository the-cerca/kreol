package controllers

import (
	"fmt"
	"html/template"
	"kreol/models"
	"kreol/statics"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LanguageController struct {
	Lm *models.LanguageManager
}

func (lc *LanguageController) HandleGetLanguages(w http.ResponseWriter, r *http.Request) {
	v := AddViews("views/selection.gohtml")
	l, err := lc.Lm.QueriesLanguages(r.Context())
	if err != nil {
		log.Printf("Error querying languages: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(statics.StaticsFs, v...)
	if err != nil {
		log.Printf("Error parsing template: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, l); err != nil {
		log.Printf("Error executing template: %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

func (lc *LanguageController) HandleSubscribeLanguage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this function was fire 🔥")
	id := chi.URLParam(r, "id")
	if id == "#" {
		return
	}
	if err := lc.Lm.SubscribeLanguage(r.Context(), id); err != nil {
		log.Printf("subscribe language : %s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	n, err := lc.Lm.QueryNameLanguageByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		log.Printf(" quering language selection : %s", err)
		return
	}
	v := fmt.Sprintf("/lang/%s", n)
	w.Header().Set("Hx-Redirect", v)
}

func (lc *LanguageController) HandleGetCourseLanguage(w http.ResponseWriter, r *http.Request) {
	lang := chi.URLParam(r, "lang")
	id, err := lc.Lm.QueryIdLanguageByName(r.Context(), lang)
	if err != nil {
		http.Error(w, "Language not found", http.StatusNotFound)
		log.Printf("QueryIdLanguageByName error: %s", err)
		return
	}
	t, err := lc.Lm.QueryAllUnseenTheme(r.Context(), id)
	if err != nil || len(*t) == 0 {
		http.Error(w, "No themes found", http.StatusInternalServerError)
		log.Printf("QueryAllUnseenTheme error: %s", err)
		return
	}
	p, err := lc.Lm.QueryWordByTheme(r.Context(), (*t)[0].ID)
	if err != nil {
		http.Error(w, "Word query failed", http.StatusInternalServerError)
		log.Printf("QueryWordByTheme error: %s", err)
		return
	}

	data := struct {
		Theme string
		Play  models.Play
	}{
		Theme: (*t)[0].Name,
		Play:  *p,
	}

	v := AddViews("views/play.gohtml")
	tpl := template.Must(template.ParseFS(statics.StaticsFs, v...))
	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Printf("Template execution error: %s", err)
	}
}
func (lc *LanguageController) PlayReponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fired")
}
