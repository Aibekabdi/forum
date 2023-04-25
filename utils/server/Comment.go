package server

import (
	"forum/utils/database"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func (db *DBase) Comment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment" {
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
	for v := range r.Form {
		if v != "comment_text" && v != "comment" {
			w.WriteHeader(400)
			ErrorPrint(w, 405)
			return
		}
	}
	text := strings.TrimFunc(r.FormValue("comment_text"), func(r rune) bool {
		return unicode.IsSpace(r)
	})
	if text == "" {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	} else {
		if err := database.InsertComment(db.Db, linkID, text, cookie); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
	}

	http.Redirect(w, r, prevURL, http.StatusFound)
	return
}

func (db *DBase) CommentRating(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ratecomment" {
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
	}
	prevURL := r.Header.Get("Referer")
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
	input := strings.Split(r.FormValue("likes"), ",")
	if len(input) != 2 {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}
	likes, err := strconv.Atoi(input[1])
	if err != nil {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}

	commentid, err := strconv.Atoi(input[0])
	if err != nil {
		w.WriteHeader(400)
		ErrorPrint(w, 400)
		return
	}

	if likes < 0 && likes > 1 {
		http.Redirect(w, r, prevURL, http.StatusFound)
	}
	if err := database.InsertCommentLike(db.Db, likes, commentid, cookie); err != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}
	http.Redirect(w, r, prevURL, http.StatusFound)
}
