package database

import (
	"database/sql"
	"errors"
	forum "forum/utils"
	"net/mail"
	"regexp"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

func IsValidRegister(db *sql.DB, user forum.User) error {
	if err := isValidEmail(user); err != nil {
		return err
	} else if err := isValidUsername(user); err != nil {
		return err
	} else if err := isValidPassword(user); err != nil {
		return err
	} else if err := NewUser(db, user); err != nil {
		return err
	}
	return nil
}

func isValidEmail(user forum.User) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}
	return nil
}

func isValidUsername(user forum.User) error {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", user.Username); !ok {
		return errors.New("invalid username")
	}
	return nil
}

func isValidPassword(user forum.User) error {
	if len(user.Password) < 8 {
		return errors.New("invalid password")
	}
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
	} {
		for _, r := range user.Password {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return errors.New("password must have at least one" + name + "character")
	}
	if user.Password != user.ConfirmPsw {
		return errors.New("Please enter the same password in both password fields")
	}
	return nil
}
