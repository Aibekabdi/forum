package server

import (
	"forum/utils/database"
	"net/http"
	"strconv"
)

func (db *DBase) RatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ratepost" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil || !database.IsInSession(db.Db, cookie) {
		if err == nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	prevURL := r.Header.Get("Referer")
	linkID, err := strconv.Atoi(prevURL[27:])
	if err != nil {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}
	check := false
	for v := range r.Form {
		if v != "likes" {
			w.WriteHeader(400)
			ErrorPrint(w, 400)
			return
		} else {
			check = true
		}
	}
	if !check {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}
	likes, err := strconv.Atoi(r.FormValue("likes"))
	if err != nil {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}
	if likes < 0 && likes > 1 {
		http.Redirect(w, r, prevURL, http.StatusFound)
	}
	if err := database.InsertPostlike(db.Db, likes, linkID, cookie); err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	http.Redirect(w, r, prevURL, http.StatusFound)

}
