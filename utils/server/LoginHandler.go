package server

import (
	"fmt"
	forum "forum/utils"
	"forum/utils/database"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

func (d *DBase) LogInHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplerr := template.ParseFiles("./templates/signin.html")
	if tmplerr != nil {
		w.WriteHeader(500)
		ErrorPrint(w, 500)
		return
	}

	if r.URL.Path != "/signin" {
		w.WriteHeader(404)
		ErrorPrint(w, 404)
		return
	}

	switch r.Method {
	// method GEt
	case "GET":
		noncorrect := forum.ErrorAuth{
			Check: false,
		}
		if _, err := r.Cookie("session"); err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if err := tmpl.Execute(w, noncorrect); err != nil {
			w.WriteHeader(500)
			ErrorPrint(w, 500)
			return
		}
	//method POST
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
			w.WriteHeader(400)
			ErrorPrint(w, 400)
			return
		}

		input := forum.User{
			Username: r.FormValue("uname"),
			Password: r.FormValue("psw"),
		}

		if err := database.IsUser(d.Db, input); err != nil {
			noncorrect := forum.ErrorAuth{
				Message: err,
				Check:   true,
			}
			tmpl.Execute(w, noncorrect)
		} else {
			id := uuid.NewV4()
			if err := database.IsLogined(d.Db, input.Username); err != nil {
				log.Println(err)
			}
			SetCookie(w, r, "session", id.String())
			if err := database.InsertCookie(d.Db, id.String(), input.Username); err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				ErrorPrint(w, 500)
				return
			}
			log.Printf("User %v is started session \n", input.Username)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	default:
		w.WriteHeader(405)
		ErrorPrint(w, 405)
		return
	}
}
