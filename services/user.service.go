package services

import (
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/repos"
)

type UserService struct {
	userRepo repos.UserRepository
}

func NewUserService(userRepo repos.UserRepository) UserService {
	return UserService{userRepo}
}

func (us *UserService) CreateUser(user *models.UserCreate) (userResponce models.UserResponse, err error) {
	userResponce, err = us.userRepo.AddUser(user)
	return
}
