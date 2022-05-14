package delivery

import (
	"html/template"
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	templates map[string]*template.Template
}

func NewHandler(service *service.Service, template map[string]*template.Template) *Handler {
	return &Handler{templates: template}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	// connecting assets
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))
	// auth
	mux.HandleFunc("/signup", h.getSignup) // signup page
	mux.HandleFunc("/post/signup", h.postSignup)
	// mux.HandleFunc("/signin", nil)         // signin page
	// mux.HandleFunc("/signout", nil)        // signout page
	// // view
	// mux.HandleFunc("/", nil)        // home page
	// mux.HandleFunc("/filter", nil)  // filter page
	// mux.HandleFunc("/profile", nil) // profile page
	// mux.HandleFunc("/post/", nil)   // post page
	// // create
	// mux.HandleFunc("/create/post", nil)    // create post page
	// mux.HandleFunc("/create/comment", nil) // create comment page
	// // rate
	// mux.HandleFunc("/rate/post", nil)    // rate post page
	// mux.HandleFunc("/rate/comment", nil) // rate comment page
	return mux
	//
}
