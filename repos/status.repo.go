package repos

import (
	"math"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type StatusRepository struct {
	DB *gorm.DB
}

func NewStatusRepository(DB *gorm.DB) StatusRepository {
	return StatusRepository{DB}
}

func (sr *StatusRepository) GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status] /* statuses []models.Status */, err error) {
	//Собираем данные для ответа в ручке с пагинацией
	sr.DB.Find(&statusesWithPaging.Items)
	statusesWithPaging.Pagination.TotalItems = len(statusesWithPaging.Items)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		statusesWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитваем количество страниц
		statusesWithPaging.Pagination.TotalPages =
			int(math.Ceil(float64(statusesWithPaging.Pagination.TotalItems) / float64(pagingOptions.PageSize)))
	}

	//Значения по умолчанию для pagingOptions
	if pagingOptions.Order != "desc" {
		pagingOptions.Order = "asc"
	}
	if pagingOptions.Page <= 0 {
		pagingOptions.Page = 1
	}
	if pagingOptions.PageSize <= 0 {
		pagingOptions.PageSize = statusesWithPaging.Pagination.TotalItems
	}
	if pagingOptions.OrderBy == "" {
		pagingOptions.OrderBy = "id"
	}
	statusesWithPaging.Pagination.Options = pagingOptions

	//Добываем выборку с учетом параметров пагинации
	err = sr.DB.Order(pagingOptions.OrderBy + " " + string(pagingOptions.Order)).
		Limit(pagingOptions.PageSize).
		Offset((pagingOptions.Page - 1) * pagingOptions.PageSize).
		Find(&statusesWithPaging.Items).Error

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

func (sr *StatusRepository) Update(updatedStatus *models.Status) (err error) {
	status := models.Status{}
	err = sr.DB.First(&status, updatedStatus.ID).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}
	status.Name = updatedStatus.Name
	err = sr.DB.Save(&status).Error
	return
}

func (sr *StatusRepository) hasTasks(statusId int) bool {
	tasksWithStatusCount := int64(0)
	sr.DB.Table("tasks").Where("status_id = ?", statusId).Count(&tasksWithStatusCount)
	return tasksWithStatusCount > 0
}

func (sr *StatusRepository) Delete(id int) (err error) {
	err = sr.DB.First(&models.Status{}, id).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}

	if sr.hasTasks(id) {
		err = common.ErrStatusHasRelatedTasks
		return
	}

	err = sr.DB.Delete(&models.Status{}, id).Error
	return
}
