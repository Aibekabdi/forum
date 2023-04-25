package forum

type User struct {
	ID         int
	Email      string
	Username   string
	Password   string
	ConfirmPsw string
}

type ErrorAuth struct {
	Message error
	Check   bool
}

type Authenticated struct {
	IsAuth bool
	Posts  []Post
}

type Post struct {
	PostID   int
	UserId   int
	Username string
	Title    string
	Tags     []string
	Content  string
	Comment  []Comment
	Like     int
	Dislike  int
}

type Comment struct {
	CommentId int
	UserId    int
	Text      string
	Username  string
	Like      int
	Dislike   int
}

type Chosen struct {
	IsAuth     bool
	ChosenPost Post
}

type Profile struct {
	UserPosts  []Post
	LikedPosts []Post
	IsAuth     bool
}

type HtmlStatus struct {
	Status string
	Text   string
}
