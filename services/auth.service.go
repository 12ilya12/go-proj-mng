package services

import (
	"errors"
	"os"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/dgrijalva/jwt-go"
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

func (as *AuthService) Login(login string, password string) (accessToken string, err error) {
	//Ищем пользователя по логину
	user, err := as.userService.FindByLogin(login)
	if err != nil {
		err = errors.New("пользователь с логином " + login + " не найден")
		return
	}

	//Сверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("неверный логин или пароль")
		return
	}

	//Создание JWT токена
	claims := &models.Claims{UserId: user.ID, Role: user.Role}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	accessToken, _ = token.SignedString([]byte(os.Getenv("token_password")))

	return
}
