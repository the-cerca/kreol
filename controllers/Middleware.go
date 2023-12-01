package controllers

import (
	"kreol/models"
	"log"
	"net/http"
)

type Middleware struct {
	Sm *models.SessionManager
}

func (m *Middleware) GetCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := ReadCookie(r, session)
		if err != nil {
			log.Println("Cookie error :", err)
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}
		u, err := m.Sm.FindUserByCookie(ctx, c.Value)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}
		ctx = models.SetUserContext(ctx, u)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
// func (m *Middleware)AuthorizedLang(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 	})
// }
