package server

import (
	forum "forum/utils"
	"forum/utils/database"
	"net/http"
	"text/template"
)

func (db *DBase) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplerr := template.ParseFiles("./templates/signup.html")
	if tmplerr != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}

	if r.URL.Path != "/signup" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}

	switch r.Method {
	// method GEt
	case "GET":
		if _, err := r.Cookie("session"); err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		empty := forum.ErrorAuth{}
		if err := tmpl.Execute(w, empty); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
	//method POST
	case "POST":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}

		input := forum.User{
			Email:      r.FormValue("email"),
			Username:   r.FormValue("uname"),
			Password:   r.FormValue("psw"),
			ConfirmPsw: r.FormValue("confirmPSW"),
		}

		if err := database.IsValidRegister(db.Db, input); err != nil {
			noncorrect := forum.ErrorAuth{
				Message: err,
				Check:   true,
			}
			tmpl.Execute(w, noncorrect)
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
	default:
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
}
