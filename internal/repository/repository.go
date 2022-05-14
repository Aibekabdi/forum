package repository

type Auth interface {
	Signin()
	Signup()
	Signout()
}

type Repository struct {
	Auth
}

func NewRepository() *Repository {
	return &Repository{}
}
