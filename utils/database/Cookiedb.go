package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func InsertCookie(db *sql.DB, uid string, login string) error {
	row := db.QueryRow("SELECT id FROM User where username = ?", login)
	var curUserId int
	err := row.Scan(&curUserId)
	if err != nil {
		return fmt.Errorf("DBuser scan error %w", err)
	}
	if _, err := db.Exec("INSERT INTO Cookie (Session, UserId) VALUES (?,?)", uid, curUserId); err != nil {
		return fmt.Errorf("inserting is invalid %w", err)
	}
	return nil
}

func IsInSession(db *sql.DB, cookie *http.Cookie) bool {
	row := db.QueryRow("SELECT Session FROM Cookie where Session = ?", cookie.Value)

	var session string
	err := row.Scan(&session)
	if err != nil {
		return false
	}
	if cookie.Value != session {
		return false
	}
	return true
}

// site  - session

func DeleteCookie(db *sql.DB, cookie *http.Cookie) {
	_, err := db.Exec("DELETE from Cookie where Session = ?", cookie.Value)
	if err != nil {
		log.Println(err)
		return
	}
}
