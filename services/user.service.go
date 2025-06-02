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

func (us *UserService) CreateUser(user *models.User) (err error) {
	//Не планируется создавать администраторов
	user.Role = "USER"

	err = us.userRepo.AddUser(user)

	//Отфильтровать какие поля возвращать в ответе
	return
}
