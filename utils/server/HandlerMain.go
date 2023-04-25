package server

import (
	"forum/utils/database"
	"net/http"
	"text/template"
)

func (d *DBase) MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
	parser := []string{
		"templates/main.html",
		"templates/header.html",
	}
	tmpl, tmplerr := template.ParseFiles(parser...)
	if tmplerr != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil || !database.IsInSession(d.Db, cookie) {
		if err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		Auth, err := database.GetPosts(d.Db)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		Auth.IsAuth = false
		if err := tmpl.Execute(w, Auth); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
		return
	}

	Auth, err := database.GetPosts(d.Db)
	if err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	Auth.IsAuth = true
	if err := tmpl.Execute(w, Auth); err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
}
