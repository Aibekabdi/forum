package server

import (
	forum "forum/utils"
	"log"
	"net/http"
	"text/template"
)

func ErrorPrint(w http.ResponseWriter, status int) {
	tmpl, tmplerr := template.ParseFiles("templates/error.html")
	if tmplerr != nil {
		w.WriteHeader(500)
		return
	}
	pr := forum.HtmlStatus{}
	if status == 404 {
		pr.Status = "Page Not Found"
		pr.Text = "We're sorry, the page you were looking for isn't found here. The link you followed may either be broken or no longer exists. Please try again, or take a look at our."
	} else if status == 500 {
		pr.Status = "Internal Server Error"
		pr.Text = "Oh no! Our web-site code is not working properly. We will be back soon!"
	} else if status == 405 {
		pr.Status = "Method Not Allowed"
		pr.Text = "Hey! Hey! Hey! Don't think that you are the most cunning)))"
	} else if status == 400 {
		pr.Status = "Bad Request"
	}
	if err := tmpl.Execute(w, pr); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}
