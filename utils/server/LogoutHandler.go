package server

import (
	"forum/utils/database"
	"net/http"
)

func (d *DBase) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signout" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	} else {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if !database.IsInSession(d.Db, cookie) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			database.DeleteCookie(d.Db, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
}
