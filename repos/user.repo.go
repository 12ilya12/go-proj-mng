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

func (ur *UserRepository) AddUser(user *models.UserCreate) (newUser models.UserResponse, err error) {
	ur.DB.Create(user)
	return
}
