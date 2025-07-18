package repos

import (
	"math"
	"strings"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(DB *gorm.DB) TaskRepository {
	return TaskRepository{DB}
}

func (sr *TaskRepository) GetAll(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task], err error) {
	//Сортировка. По умолчанию по возрастанию идентификатора.
	var orderRule string
	if pagingOptions.OrderBy == "" {
		orderRule = "id"
	} else {
		var columnCount int64
		sr.DB.Select("column_name").Table("information_schema.columns").
			Where("table_name = ? AND column_name = ?", "tasks", pagingOptions.OrderBy).Count(&columnCount)
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

	//Фильтрация
	if taskFilters.StatusId > 0 {
		tx = tx.Where("status_id = ?", taskFilters.StatusId)
	}
	if taskFilters.CategoryId > 0 {
		tx = tx.Where("category_id = ?", taskFilters.CategoryId)
	}

	if strings.ToLower(taskFilters.UserRole) == "admin" {
		err = tx.Find(&tasksWithPaging.Items).Error
	} else {
		//Обычному пользователю показываем только его задачи
		err = tx.Find(&tasksWithPaging.Items, "user_id = ?", taskFilters.UserId).Error
	}

	//Собираем выходные данные пагинации
	txForTotalItems := sr.DB.Model(&models.Task{}) //Подсчёт количества задач с учетом фильтраций
	if taskFilters.StatusId > 0 {
		txForTotalItems = txForTotalItems.Where("status_id = ?", taskFilters.StatusId)
	}
	if taskFilters.CategoryId > 0 {
		txForTotalItems = txForTotalItems.Where("category_id = ?", taskFilters.CategoryId)
	}
	if strings.ToLower(taskFilters.UserRole) == "admin" {
		txForTotalItems.Count(&tasksWithPaging.Pagination.TotalItems)
	} else {
		txForTotalItems.Where("user_id = ?", taskFilters.UserId).Count(&tasksWithPaging.Pagination.TotalItems)
	}
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		tasksWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитываем количество страниц
		tasksWithPaging.Pagination.TotalPages =
			int64(math.Ceil(float64(tasksWithPaging.Pagination.TotalItems) /
				float64(pagingOptions.PageSize)))
	}
	tasksWithPaging.Pagination.Options = pagingOptions

	return
}

func (sr *TaskRepository) GetById(id uint) (task models.Task, err error) {
	err = sr.DB.First(&task, id).Error
	return
}

func (sr *TaskRepository) Create(task *models.Task) (err error) {
	err = sr.DB.Create(&task).Error
	return
}

func (sr *TaskRepository) Update(updatedTask *models.Task) (err error) {
	err = sr.DB.Save(&updatedTask).Error
	return
}

func (sr *TaskRepository) HasDependencies(taskId uint) bool {
	var depsWithTaskCount int64
	sr.DB.Table("dependencies").Where("parent_task_id = ? OR child_task_id = ?", taskId, taskId).Count(&depsWithTaskCount)
	return depsWithTaskCount > 0
}

func (sr *TaskRepository) Delete(id uint) (err error) {
	err = sr.DB.Delete(&models.Task{}, id).Error
	return
}
