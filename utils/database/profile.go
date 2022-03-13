package database

import (
	"database/sql"
	forum "forum/utils"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func FilterProfile(db *sql.DB, cookie *http.Cookie) (forum.Profile, error) {
	myprofile := forum.Profile{}
	id, _, err := GetUser(db, cookie)
	if err != nil {
		return myprofile, err
	}
	rows, err := db.Query("SELECT id, PostTitle, Tags , PostText FROM Post where UserId = ?", id)
	if err != nil {
		log.Println(err)
		return myprofile, err
	}
	defer rows.Close()
	mypost := []forum.Post{}
	for rows.Next() {
		p := forum.Post{}
		var stringtag string
		err := rows.Scan(&p.PostID, &p.Title, &stringtag, &p.Content)
		p.Tags = strings.Split(stringtag, ",")
		if err != nil {
			return myprofile, err
		}
		mypost = append(mypost, p)
	}

	likedpost := []forum.Post{}
	newrows, err := db.Query("SELECT Post.id, Post.PostTitle, Post.Tags , Post.PostText, Post.UserId FROM Post JOIN PostRating on PostRating.Postid = Post.id where LikedUserid = ? and likes = 0", id)
	if err != nil {
		log.Println(err, "newrows")
		return myprofile, err
	}
	defer newrows.Close()
	for newrows.Next() {
		p := forum.Post{}
		var stringtag string
		err := newrows.Scan(&p.PostID, &p.Title, &stringtag, &p.Content, &p.UserId)
		row := db.QueryRow("SELECT Username FROM User where id = ?", p.UserId)
		if eerr := row.Scan(&p.Username); eerr != nil {
			log.Println(eerr)
			return myprofile, eerr
		}
		p.Tags = strings.Split(stringtag, ",")
		if err != nil {
			return myprofile, err
		}
		likedpost = append(likedpost, p)
	}

	myprofile = forum.Profile{
		IsAuth:     true,
		UserPosts:  mypost,
		LikedPosts: likedpost,
	}

	return myprofile, nil
}
