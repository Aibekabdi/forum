package database

import (
	"database/sql"
	"errors"
	"fmt"
	forum "forum/utils"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func NewUser(db *sql.DB, user forum.User) error {
	hash, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("hashing error: %w", err)
	}
	if _, err := db.Exec("INSERT INTO User (email, username, password) VALUES (?,?,?)", user.Email, user.Username, hash); err != nil {
		newerr := errors.New("email or username is already exist, try another ones")
		return newerr
	}
	return nil
}

func InsertPostlike(db *sql.DB, likes int, postid int, cookie *http.Cookie) error {
	userid, _, err := GetUser(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	update := "UPDATE PostRating SET Likes = ? WHERE Postid = ? AND LikedUserid = ?"
	delete := "DELETE FROM PostRating WHERE Postid = ? AND LikedUserid = ?"
	var curlikes int
	row := db.QueryRow("SELECT Likes FROM PostRating where Postid = ? and LikedUserid = ?", postid, userid)
	ers := row.Scan(&curlikes)
	if ers == nil {
		if (curlikes == 0 && likes == 0) || (curlikes == 1 && likes == 1) {
			if _, err = db.Exec(delete, postid, userid); err != nil {
				return err
			}
		} else if (curlikes == 0 && likes == 1) || (curlikes == 1 && likes == 0) {
			if _, err = db.Exec(update, likes, postid, userid); err != nil {
				return err
			}
		}
	} else {
		if _, err := db.Exec("INSERT INTO PostRating (Likes, Postid, LikedUserid) VALUES (?,?,?)", likes, postid, userid); err != nil {
			log.Println(ers)
			return err
		}
	}
	return nil
}

func InsertComment(db *sql.DB, postid int, text string, cookie *http.Cookie) error {
	userid, _, err := GetUser(db, cookie)
	if err != nil {
		return err
	}
	if _, err := db.Exec("INSERT INTO Comments(CommenterId, PostId, CommentText) VALUES(?,?,?) ", userid, postid, text); err != nil {
		return err
	}
	return nil
}

func InsertCommentLike(db *sql.DB, likes int, commentid int, cookie *http.Cookie) error {
	userid, _, err := GetUser(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	update := "UPDATE CommentRating SET Likes = ? WHERE Commentid = ? AND LikedUserid = ?"
	delete := "DELETE FROM CommentRating WHERE Commentid = ? AND LikedUserid = ?"
	var curlikes int
	row := db.QueryRow("SELECT Likes FROM CommentRating where Commentid = ? and LikedUserid = ?", commentid, userid)
	ers := row.Scan(&curlikes)
	if ers == nil {
		if (curlikes == 0 && likes == 0) || (curlikes == 1 && likes == 1) {
			if _, err = db.Exec(delete, commentid, userid); err != nil {
				return err
			}
		} else if (curlikes == 0 && likes == 1) || (curlikes == 1 && likes == 0) {
			if _, err = db.Exec(update, likes, commentid, userid); err != nil {
				return err
			}
		}
	} else {
		if _, err := db.Exec("INSERT INTO CommentRating (Likes, Commentid, LikedUserid) VALUES (?,?,?)", likes, commentid, userid); err != nil {
			log.Println(ers)
			log.Println(err)
			return err
		}
	}
	return nil
}
