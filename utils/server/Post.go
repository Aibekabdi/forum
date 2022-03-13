package server

import (
	forum "forum/utils"
	"forum/utils/database"
	"log"
	"net/http"
	"text/template"
)

func (db *DBase) CurrentPost(w http.ResponseWriter, r *http.Request) {
	parser := []string{
		"templates/post.html",
		"templates/header.html",
	}
	tmpl, tmplerr := template.ParseFiles(parser...)
	if tmplerr != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
	id := r.URL.Path[6:]
	check, postid := database.IsPost(db.Db, id)
	if !check {
		ErrorPrint(w, 404)
		ErrorPrint(w, 404)
		return
	}
	post, err := database.ChosenPost(db.Db, postid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	Auth := forum.Chosen{
		ChosenPost: post,
	}
	cookie, err := r.Cookie("session")
	if err != nil || !database.IsInSession(db.Db, cookie) {
		if err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		Auth.IsAuth = false
	} else {
		Auth.IsAuth = true
	}
	if err := tmpl.Execute(w, Auth); err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
}
