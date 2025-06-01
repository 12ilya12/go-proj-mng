package services

import "github.com/12ilya12/go-proj-mng/models"

type AuthService struct {
	userService UserService
}

// Auth UseCase constructor
func NewAuthService(userService UserService) AuthService {
	return AuthService{userService}
}

func (as *AuthService) Register(user *models.UserCreate) (userResponce models.UserResponse, err error) {
	userResponce, err = as.userService.CreateUser(user)
	return
}
