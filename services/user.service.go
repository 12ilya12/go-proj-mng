package services

import (
	"github.com/12ilya12/go-proj-mng/repos"
)

type UserService struct {
	userRepo repos.UserRepository
}

func NewUserService(userRepo repos.UserRepository) UserService {
	return UserService{userRepo}
}
