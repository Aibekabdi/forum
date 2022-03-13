package database

import "database/sql"

func CreateDB(db *sql.DB) error {
	//Create tables
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS User(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL, 
		username TEXT UNIQUE NOT NULL, 
		password TEXT NOT NULL
	)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS Cookie(
		Session Text NOT NULL,
		UserId INTEGER NOT NULL,
		FOREIGN KEY(UserId) REFERENCES User(id) ON DELETE CASCADE
	)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS Post(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		UserId INTEGER NOT NULL,
		PostTitle Text NOT NULL,
		Tags Text,
		PostText Text NOT NULL,
		FOREIGN KEY(UserId) REFERENCES User(id) ON DELETE CASCADE
	)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS Comments(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		CommenterId INTEGER NOT NULL,
		PostId INTEGER NOT NULL,
		CommentText text NOT NULL, 
		FOREIGN KEY(CommenterId) REFERENCES User(id) ON DELETE CASCADE,
		FOREIGN KEY(PostId) REFERENCES Post(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS PostRating(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Likes int,
		Postid INTEGER NOT NULL,
		LikedUserid INTEGER NOT NULL,
		FOREIGN KEY(Postid) REFERENCES Post(id) ON DELETE CASCADE,
		FOREIGN KEY(LikedUserid) REFERENCES User(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS CommentRating(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Likes int,
		Commentid INTEGER NOT NULL,
		LikedUserid INTEGER NOT NULL,
		FOREIGN KEY(Commentid) REFERENCES Comments(id) ON DELETE CASCADE,
		FOREIGN KEY(LikedUserid) REFERENCES User(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}
	return nil
}