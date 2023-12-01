package controllers

import (
	"fmt"
	"net/http"
)

const (
	session = "session_id"
)

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)
}
func ReadCookie(r *http.Request, name string) (*http.Cookie, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("no cookie with the name %s", name)
	}
	return c, nil
}
