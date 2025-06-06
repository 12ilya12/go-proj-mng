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

func (sr *StatusRepository) Create(status *models.Status) (err error) {
	err = sr.DB.Create(&status).Error
	return
}

func (sr *StatusRepository) Update(id int, newStatus *models.Status) (err error) {
	status := models.Status{}
	err = sr.DB.First(&status, id).Error
	if err != nil {
		//TODO: Вернуть ошибку со статусом NotFound
		return
	}
	status.Name = newStatus.Name
	err = sr.DB.Save(&status).Error
	if err != nil {
		//TODO: Тут тоже нужно записать статус ошибки.
	}
	return
}

func (sr *StatusRepository) Delete(id int) (err error) {
	err = sr.DB.Delete(&models.Status{}, id).Error
	return
}
