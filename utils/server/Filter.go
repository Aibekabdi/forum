package server

import (
	"forum/utils/database"
	"log"
	"net/http"
	"text/template"
)

func (db *DBase) FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		ErrorPrint(w, 404)
		w.WriteHeader(404)
		return
	}

	if r.Method == "GET" {
		ErrorPrint(w, 405)
		w.WriteHeader(405)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if r.Method != "POST" {
		ErrorPrint(w, 405)
		w.WriteHeader(405)
		return
	}

	parser := []string{
		"templates/main.html",
		"templates/header.html",
	}

	tmpl, tmplerr := template.ParseFiles(parser...)
	if tmplerr != nil {
		ErrorPrint(w, 500)
		w.WriteHeader(500)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		ErrorPrint(w, 400)
		w.WriteHeader(400)
		return
	}

	input := r.FormValue("filter")
	m := database.MapTag()
	if _, ok := m[input]; !ok {
		ErrorPrint(w, 400)
		w.WriteHeader(400)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || !database.IsInSession(db.Db, cookie) {
		if err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}

		Auth, err := database.GetChosenPosts(db.Db, input)

		if err != nil {
			ErrorPrint(w, 500)
			w.WriteHeader(500)
			return
		}

		Auth.IsAuth = false

		if err := tmpl.Execute(w, Auth); err != nil {
			ErrorPrint(w, 500)
			w.WriteHeader(500)
			return
		}
		return
	}

	Auth, err := database.GetChosenPosts(db.Db, input)
	if err != nil {
		ErrorPrint(w, 500)
		w.WriteHeader(500)
		return
	}

	Auth.IsAuth = true

	if err := tmpl.Execute(w, Auth); err != nil {
		ErrorPrint(w, 500)
		w.WriteHeader(500)
		return
	}
}
