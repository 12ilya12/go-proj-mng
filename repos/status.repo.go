package repos

import (
	"github.com/12ilya12/go-proj-mng/models"
	"gorm.io/gorm"
)

type StatusRepository struct {
	DB *gorm.DB
}

func NewStatusRepository(DB *gorm.DB) StatusRepository {
	return StatusRepository{DB}
}

func (sr *StatusRepository) GetAll() (statuses []models.Status, err error) {
	err = sr.DB.Find(&statuses).Error
	return
}

func (sr *StatusRepository) GetById(id int) (status models.Status, err error) {
	err = sr.DB.First(&status, id).Error
	return
}
