package services

import (
	"github.com/12ilya12/go-proj-mng/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService UserService
}

// Auth UseCase constructor
func NewAuthService(userService UserService) AuthService {
	return AuthService{userService}
}

func (as *AuthService) Register(user *models.User) (err error) {
	//Хэшируем пароль
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	err = as.userService.CreateUser(user)

	//Тут можно создать токен пользователя...

	//Удаляем пароль
	user.Password = ""

	return
}
