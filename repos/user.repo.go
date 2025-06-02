package repos

import (
	"errors"

	"github.com/12ilya12/go-proj-mng/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return UserRepository{DB}
}

func (ur *UserRepository) AddUser(user *models.User) (err error) {

	ur.DB.Create(user)

	if user.ID <= 0 {
		err = errors.New("Ошибка при создании пользователя в базе данных")
	}

	return
}
