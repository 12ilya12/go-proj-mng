package repos

import (
	"math"
	"strings"

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
