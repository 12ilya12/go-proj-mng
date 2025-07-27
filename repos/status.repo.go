package repos

import (
	"math"
	"strings"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type StatusRepository interface {
	GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status], err error)
	GetById(id uint) (status models.Status, err error)
	Create(status *models.Status) (err error)
	Update(paramsForUpdate *models.Status) (updatedStatus models.Status, err error)
	HasTasks(statusId uint) bool
	Delete(id uint) (err error)
}

type StatusRepositoryImpl struct {
	DB *gorm.DB
}

func NewStatusRepositoryImpl(DB *gorm.DB) StatusRepository {
	return &StatusRepositoryImpl{DB}
}

func (sr *StatusRepositoryImpl) GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status] /* statuses []models.Status */, err error) {
	//Сортировка. По умолчанию по возрастанию идентификатора.
	var orderRule string
	if pagingOptions.OrderBy == "" {
		orderRule = "id"
	} else {
		var columnCount int64
		sr.DB.Select("column_name").Table("information_schema.columns").
			Where("table_name = ? AND column_name = ?", "statuses", pagingOptions.OrderBy).Count(&columnCount)
		if columnCount == 0 {
			//Колонка, по которой планировалось сортировать, отсутствует в таблице
			pagingOptions.OrderBy = ""
			orderRule = "id"
		} else {
			orderRule = pagingOptions.OrderBy
		}
	}
	if strings.ToLower(string(pagingOptions.Order)) == "desc" {
		orderRule += " desc"
	}
	tx := sr.DB.Order(orderRule)

	//Пагинация
	if pagingOptions.PageSize > 0 {
		tx = tx.Limit(pagingOptions.PageSize)
	}
	if pagingOptions.Page > 0 {
		tx = tx.Offset((pagingOptions.Page - 1) * pagingOptions.PageSize)
	}

	err = tx.Find(&statusesWithPaging.Items).Error

	//Собираем выходные данные пагинации
	sr.DB.Model(&models.Status{}).Count(&statusesWithPaging.Pagination.TotalItems)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		statusesWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитываем количество страниц
		statusesWithPaging.Pagination.TotalPages =
			int64(math.Ceil(float64(statusesWithPaging.Pagination.TotalItems) /
				float64(pagingOptions.PageSize)))
	}
	statusesWithPaging.Pagination.Options = pagingOptions

	return
}

func (sr *StatusRepositoryImpl) GetById(id uint) (status models.Status, err error) {
	err = sr.DB.First(&status, id).Error
	return
}

func (sr *StatusRepositoryImpl) Create(status *models.Status) (err error) {
	err = sr.DB.Create(&status).Error
	return
}

func (sr *StatusRepositoryImpl) Update(paramsForUpdate *models.Status) (updatedStatus models.Status, err error) {
	updatedStatus = models.Status{}
	err = sr.DB.First(&updatedStatus, paramsForUpdate.ID).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}
	updatedStatus.Name = paramsForUpdate.Name
	err = sr.DB.Save(&updatedStatus).Error
	return
}

func (sr *StatusRepositoryImpl) HasTasks(statusId uint) bool {
	var tasksWithStatusCount int64
	sr.DB.Table("tasks").Where("status_id = ?", statusId).Count(&tasksWithStatusCount)
	return tasksWithStatusCount > 0
}

func (sr *StatusRepositoryImpl) Delete(id uint) (err error) {
	err = sr.DB.Delete(&models.Status{}, id).Error
	return
}
