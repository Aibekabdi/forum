package server

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, r *http.Request, name string, value string) {
	cookie, err := r.Cookie(name)
	time := time.Now().Add(30 * time.Minute)
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:    name,
			Value:   value,
			Expires: time,
			MaxAge:  3600,
		}
	}
	http.SetCookie(w, cookie)
}
