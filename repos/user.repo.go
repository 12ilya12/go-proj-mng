package repos

import (
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
	err = ur.DB.Create(user).Error
	return
}

func (ur *UserRepository) FindByLogin(login string) (user models.User, err error) {
	err = ur.DB.First(&user, "login = ?", login).Error
	return
}
