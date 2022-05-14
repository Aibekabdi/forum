package delivery

import (
	"log"
	"net/http"
)

func (h *Handler) getSignup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		w.WriteHeader(404)
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	tmpl := h.templates["signup.html"]
	if tmpl == nil {
		// could not find it
		w.WriteHeader(500)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (h *Handler) postSignup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/signup" {
		w.WriteHeader(404)
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	log.Println("hello world")
}
