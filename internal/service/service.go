package service

import "forum/internal/repository"

type Auth interface {
	Signup()
	Signin()
	Signout()
}

type Service struct{}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
