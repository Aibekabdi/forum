package database

import (
	"database/sql"
	"errors"
	forum "forum/utils"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)
 
func IsUser(db *sql.DB, user forum.User) error {
	row := db.QueryRow("SELECT username , password  FROM User where username = ?", user.Username)
	curUser := forum.User{}
	err := row.Scan(&curUser.Username, &curUser.Password)
	if !CheckPasswordHash(user.Password, curUser.Password) {
		return errors.New("non correct login or password")
	} else if err != nil {
		return errors.New("non correct login or password")
	}
	return nil
}

func GetUser(db *sql.DB, c *http.Cookie) (int, string, error) {
	row := db.QueryRow("SELECT User.id, User.username  FROM User JOIN Cookie on UserId = id where Session = ?", c.Value)
	var UserId int
	var Uname string
	err := row.Scan(&UserId, &Uname)
	if err != nil {
		return 0, "", err
	}
	return UserId, Uname, nil
}

func IsLogined(db *sql.DB, username string) error {
	row := db.QueryRow("SELECT id FROM User where username = ?", username)
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE from Cookie where UserId = ?", id); err != nil {
		return err
	}
	return nil
}
