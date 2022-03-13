package main

import (
	"fmt"
	"forum/utils/database"
	"forum/utils/server"
	"log"
	"net/http"
)

func main() {
	//Database connecting
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.CreateDB(db); err != nil {
		log.Fatal(err)
	}
	//Database
	dbase := &server.DBase{
		Db: db,
	}
	//css connecting
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	//web sites
	mux.HandleFunc("/", dbase.MainHandler)
	mux.HandleFunc("/signup", dbase.SignUpHandler)
	mux.HandleFunc("/signin", dbase.LogInHandler)
	mux.HandleFunc("/signout", dbase.LogOut)

	mux.HandleFunc("/CreatePost", dbase.CreatePost)
	mux.HandleFunc("/post/", dbase.CurrentPost)
	mux.HandleFunc("/ratepost", dbase.RatePost)

	mux.HandleFunc("/comment", dbase.Comment)
	mux.HandleFunc("/ratecomment", dbase.CommentRating)
	mux.HandleFunc("/filter", dbase.FilterHandler)
	mux.HandleFunc("/Profile", dbase.ProfileHandler)
	//Server is listenning
	fmt.Println("http://localhost:8080/")
	fmt.Println("Server is Listenning...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
