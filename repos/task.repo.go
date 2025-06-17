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

func (sr *TaskRepository) Update(paramsForUpdate *models.Task, userInfo common.UserInfo) (updatedTask models.Task, err error) {
	updatedTask = models.Task{}
	err = sr.DB.First(&updatedTask, paramsForUpdate.ID).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}
	//Разные возможности в зависимости от роли пользователя
	if strings.ToLower(userInfo.UserRole) == "admin" {
		//Администратор может изменять все поля задачи
		if paramsForUpdate.Name != "" {
			updatedTask.Name = paramsForUpdate.Name
		}
		if paramsForUpdate.Description != "" {
			updatedTask.Description = paramsForUpdate.Description
		}
		if paramsForUpdate.StatusId != 0 {
			updatedTask.StatusId = paramsForUpdate.StatusId
		}
		if paramsForUpdate.CategoryId != 0 {
			updatedTask.CategoryId = paramsForUpdate.CategoryId
		}
		if paramsForUpdate.UserId != 0 {
			updatedTask.UserId = paramsForUpdate.UserId
		}
		if !paramsForUpdate.Deadline.IsZero() {
			updatedTask.Deadline = paramsForUpdate.Deadline
		}
		if paramsForUpdate.Priority != 0 {
			updatedTask.Priority = paramsForUpdate.Priority
		}
	} else {
		//Пользователь не администратор, поэтому он может менять только статус СВОЕЙ задачи
		if updatedTask.UserId == uint(userInfo.UserId) {
			if paramsForUpdate.StatusId != 0 {
				updatedTask.StatusId = paramsForUpdate.StatusId
			}
		} else {
			err = common.ErrUserHasNotPermissionToEditTask
			return
		}
	}
	err = sr.DB.Save(&updatedTask).Error
	return
}

func (sr *TaskRepository) hasDependencies(taskId uint) bool {
	var depsWithTaskCount int64
	sr.DB.Table("dependencies").Where("parent_task_id = ? OR child_task_id = ?", taskId, taskId).Count(&depsWithTaskCount)
	return depsWithTaskCount > 0
}

func (sr *TaskRepository) Delete(id uint) (err error) {
	err = sr.DB.First(&models.Task{}, id).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}

	if sr.hasDependencies(id) {
		err = common.ErrTaskHasRelatedDependency
		return
	}

	err = sr.DB.Delete(&models.Task{}, id).Error
	return
}
