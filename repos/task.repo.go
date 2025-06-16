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
		err = tx.Find(&tasksWithPaging.Items).Where("user_id = ?", taskFilters.UserId).Error
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
		txForTotalItems.Count(&tasksWithPaging.Pagination.TotalItems).Where("user_id = ?", taskFilters.UserId)
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

func (sr *TaskRepository) GetById(id int) (task models.Task, err error) {
	err = sr.DB.First(&task, id).Error
	return
}

func (sr *TaskRepository) Create(task *models.Task) (err error) {
	err = sr.DB.Create(&task).Error
	return
}

func (sr *TaskRepository) Update(updatedTask *models.Task, userInfo common.UserInfo) (err error) {
	task := models.Task{}
	err = sr.DB.First(&task, updatedTask.ID).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}
	//Разные возможности в зависимости от роли пользователя
	if strings.ToLower(userInfo.UserRole) == "admin" {
		//Администратор может изменять все поля задачи
		if updatedTask.Name != "" {
			task.Name = updatedTask.Name
		}
		if updatedTask.Description != "" {
			task.Description = updatedTask.Description
		}
		if updatedTask.StatusId != 0 {
			task.StatusId = updatedTask.StatusId
		}
		if updatedTask.CategoryId != 0 {
			task.CategoryId = updatedTask.CategoryId
		}
		if updatedTask.UserId != 0 {
			task.UserId = updatedTask.UserId
		}
		if !updatedTask.Deadline.IsZero() {
			task.Deadline = updatedTask.Deadline
		}
		if updatedTask.Priority != 0 {
			task.Priority = updatedTask.Priority
		}
	} else {
		//Пользователь не администратор, поэтому он может менять только статус СВОЕЙ задачи
		if task.UserId == uint32(userInfo.UserId) {
			if updatedTask.StatusId != 0 {
				task.StatusId = updatedTask.StatusId
			}
		} else {
			err = common.ErrUserHasNotPermissionToEditTask
			return
		}
	}
	err = sr.DB.Save(&task).Error
	return
}

func (sr *TaskRepository) hasTasks(taskId int) bool {
	tasksWithTaskCount := int64(0)
	sr.DB.Table("tasks").Where("task_id = ?", taskId).Count(&tasksWithTaskCount)
	return tasksWithTaskCount > 0
}

func (sr *TaskRepository) Delete(id int) (err error) {
	err = sr.DB.First(&models.Task{}, id).Error
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}

	if sr.hasTasks(id) {
		err = common.ErrTaskHasRelatedDependency
		return
	}

	err = sr.DB.Delete(&models.Task{}, id).Error
	return
}
