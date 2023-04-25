package database

import (
	"database/sql"
	"errors"
	forum "forum/utils"
	"log"
	"strconv"
	"strings"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(db *sql.DB, post forum.Post) error {
	if _, err := db.Exec("INSERT INTO Post (UserId, PostTitle, Tags, PostText) VALUES (?,?,?,?)", post.UserId, post.Title, strings.Join(post.Tags, ","), post.Content); err != nil {
		return err
	}
	return nil
}

func GetPosts(db *sql.DB) (forum.Authenticated, error) {
	auth := forum.Authenticated{}
	rows, err := db.Query("SELECT * from Post")
	if err != nil {
		return auth, err
	}
	post := []forum.Post{}
	defer rows.Close()
	for rows.Next() {
		p := forum.Post{}
		var stringtag string
		err := rows.Scan(&p.PostID, &p.UserId, &p.Title, &stringtag, &p.Content)
		p.Tags = strings.Split(stringtag, ",")
		row := db.QueryRow("SELECT Username FROM User where id = ?", p.UserId)
		if eerr := row.Scan(&p.Username); eerr != nil {
			log.Println(eerr)
			return auth, eerr
		}
		if err != nil {
			log.Println(err)
			return auth, err
		}
		post = append(post, p)
	}
	auth.Posts = post
	return auth, nil
}

func GetChosenPosts(db *sql.DB, tag string) (forum.Authenticated, error) {
	auth := forum.Authenticated{}
	rows, err := db.Query("SELECT * from Post")
	if err != nil {
		return auth, err
	}
	defer rows.Close()
	post := []forum.Post{}
	for rows.Next() {
		p := forum.Post{}
		var stringtag string
		err := rows.Scan(&p.PostID, &p.UserId, &p.Title, &stringtag, &p.Content)
		p.Tags = strings.Split(stringtag, ",")
		var check bool
		for _, v := range p.Tags {
			if v == tag {
				check = true
				break
			}
		}
		if !check {
			continue
		}

		row := db.QueryRow("SELECT Username FROM User where id = ?", p.UserId)
		if eerr := row.Scan(&p.Username); eerr != nil {
			log.Println(eerr)
			return auth, eerr
		}
		if err != nil {
			log.Println(err)
			return auth, err
		}
		post = append(post, p)
	}
	auth.Posts = post
	return auth, nil
}

func IsPost(db *sql.DB, id string) (bool, int) {
	nmb, err := strconv.Atoi(id)
	if err != nil {
		return false, 0
	}
	var curid int
	err = db.QueryRow("SELECT id FROM Post WHERE id = ?", nmb).Scan(&curid)
	if err == nil && err == sql.ErrNoRows {
		return false, nmb
	} else if curid == 0 {
		return false, 0
	}
	return true, nmb
}

func ChosenPost(db *sql.DB, postid int) (forum.Post, error) {
	var post = forum.Post{}
	var stringtag string
	row := db.QueryRow("SELECT * FROM Post where id = ?", postid)
	err := row.Scan(&post.PostID, &post.UserId, &post.Title, &stringtag, &post.Content)
	if err != nil {
		return post, err
	}
	post.Username, err = PostUserName(db, post.UserId)
	if err != nil {
		return post, err
	}

	post.Tags = strings.Split(stringtag, ",")
	post.Like, post.Dislike, err = ChosenPostlikes(db, postid)
	if err != nil {
		return post, err
	}
	post.Comment, err = ChosenPostComments(db, postid)
	if err != nil {
		return post, err
	}
	return post, nil
}

func ChosenPostlikes(db *sql.DB, postid int) (int, int, error) {
	var likes int
	var dislikes int
	rows, err := db.Query("SELECT Likes FROM PostRating where Postid =?", postid)
	if err != nil {
		return 0, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var islike int
		err := rows.Scan(&islike)
		if err != nil {
			return 0, 0, err
		}
		if islike == 0 {
			likes++
		} else if islike == 1 {
			dislikes++
		}
	}
	return likes, dislikes, nil
}

func ChosenPostComments(db *sql.DB, postid int) ([]forum.Comment, error) {
	var comments []forum.Comment
	rows, err := db.Query("SELECT id, CommenterId, CommentText FROM Comments where PostId = ?", postid)
	if err != nil {
		return comments, err
	}

	defer rows.Close()
	for rows.Next() {
		var comment forum.Comment
		err := rows.Scan(&comment.CommentId, &comment.UserId, &comment.Text)
		if err != nil {
			return comments, err
		}
		row := db.QueryRow("SELECT username FROM User where id =?", comment.UserId)
		if err = row.Scan(&comment.Username); err != nil {
			return comments, err
		}
		comment.Like, comment.Dislike, err = CommentsLikes(db, comment.CommentId)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func CommentsLikes(db *sql.DB, commentid int) (int, int, error) {
	var likes int
	var dislikes int
	rows, err := db.Query("SELECT Likes FROM CommentRating where Commentid =?", commentid)
	if err != nil {
		return 0, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var islike int
		err := rows.Scan(&islike)
		if err != nil {
			return 0, 0, err
		}
		if islike == 0 {
			likes++
		} else if islike == 1 {
			dislikes++
		}
	}
	return likes, dislikes, nil
}

func PostUserName(db *sql.DB, UserId int) (string, error) {
	row := db.QueryRow("SELECT username FROM User where id = ?", UserId)
	var Uname string
	err := row.Scan(&Uname)
	if err != nil {
		return "", err
	}
	return Uname, nil
}

func IsValidPost(post forum.Post) error {
	m := MapTag()
	tag := post.Tags
	title := strings.TrimFunc(post.Title, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	if title == "" || len(title) > 60 {
		return errors.New("invalid title")
	}
	text := strings.TrimFunc(post.Content, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	if text == "" || len(text) > 500 {
		return errors.New("invalid text")
	}
	if len(tag) == 0 {
		return errors.New("no chosen tag")
	}
	for _, v := range tag {
		if _, ok := m[v]; !ok {
			var newtag []string
			for _, v1 := range tag {
				if v != v1 {
					newtag = append(newtag, v1)
				}
			}
			if len(newtag) == 0 {
				return errors.New("no chosen tag")
			}
			post.Tags = newtag
		}
	}
	return nil
}

func MapTag() map[string]bool {
	m := make(map[string]bool)
	m["anime"] = true
	m["alem"] = true
	m["golang"] = true
	m["forum"] = true
	m["memes"] = true
	m["community"] = true
	m["other"] = true
	return m
}
