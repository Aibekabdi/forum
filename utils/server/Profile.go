package server

import (
	"forum/utils/database"
	"log"
	"net/http"
	"text/template"
)

func (db *DBase) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Profile" {
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
		"templates/profile.html",
		"templates/header.html",
	}

	tmpl, tmplerr := template.ParseFiles(parser...)
	if tmplerr != nil {
		log.Println(tmplerr)
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || !database.IsInSession(db.Db, cookie) {
		if err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	Auth, err := database.FilterProfile(db.Db, cookie)
	if err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 405)
		return
	}

	Auth.IsAuth = true
	if err := tmpl.Execute(w, Auth); err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 405)
		return
	}
}
