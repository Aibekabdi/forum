package server

import (
	forum "forum/utils"
	"forum/utils/database"
	"log"
	"net/http"
	"text/template"
)

func (db *DBase) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/CreatePost" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}
	tmpl, tmplerr := template.ParseFiles("templates/CreatePost.html")
	if tmplerr != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("session")
		if err != nil || !database.IsInSession(db.Db, cookie) {
			if err == nil {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		auth := forum.ErrorAuth{
			Check: false,
		}
		if err := tmpl.Execute(w, auth); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(400)
			ErrorPrint(w, 400)
			return
		}
		cookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(400)
			ErrorPrint(w, 400)
			return
		}
		UserID, _, err := database.GetUser(db.Db, cookie)
		if err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
		post := forum.Post{
			UserId:  UserID,
			Title:   r.FormValue("title"),
			Content: r.FormValue("postContent"),
			Tags:    r.Form["Tags"],
		}
		if err := database.IsValidPost(post); err != nil {
			auth := forum.ErrorAuth{
				Message: err,
				Check:   true,
			}
			w.WriteHeader(400)
			if terr := tmpl.Execute(w, auth); terr != nil {
				w.WriteHeader(500)
				ErrorPrint(w, 500)
				return
			}
			return
		}
		if err := database.CreatePost(db.Db, post); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)

	default:
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
}
