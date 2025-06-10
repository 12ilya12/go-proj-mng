package repos

import (
	"math"

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

// TODO: Произвести рефакторинг функций GetAll и GetAllForUser (возможно объединить в одну). Много дублирования кода и не учитывается, что некоторые параметры могут быть не заданы
func (sr *TaskRepository) GetAll(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task] /* taskes []models.Task */, err error) {
	//Собираем данные для ответа в ручке с пагинацией
	sr.DB.Find(&tasksWithPaging.Items)
	tasksWithPaging.Pagination.TotalItems = len(tasksWithPaging.Items)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		tasksWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитваем количество страниц
		tasksWithPaging.Pagination.TotalPages =
			int(math.Ceil(float64(tasksWithPaging.Pagination.TotalItems) / float64(pagingOptions.PageSize)))
	}

	//Значения по умолчанию для pagingOptions
	if pagingOptions.Order != "desc" {
		pagingOptions.Order = "asc"
	}
	if pagingOptions.Page <= 0 {
		pagingOptions.Page = 1
	}
	if pagingOptions.PageSize <= 0 {
		pagingOptions.PageSize = tasksWithPaging.Pagination.TotalItems
	}
	if pagingOptions.OrderBy == "" {
		pagingOptions.OrderBy = "id"
	}
	tasksWithPaging.Pagination.Options = pagingOptions

	//Добываем выборку с учетом параметров пагинации
	err = sr.DB.Order(pagingOptions.OrderBy+" "+string(pagingOptions.Order)).
		Limit(pagingOptions.PageSize).
		Offset((pagingOptions.Page-1)*pagingOptions.PageSize).
		Find(&tasksWithPaging.Items).Where("status_id = ?, category_id = ?", taskFilters.StatusId, taskFilters.CategoryId).Error

	return
}

func (sr *TaskRepository) GetAllForUser(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task] /* taskes []models.Task */, err error) {
	//Собираем данные для ответа в ручке с пагинацией
	sr.DB.Find(&tasksWithPaging.Items).Where("user_id = ?", taskFilters.UserId)
	tasksWithPaging.Pagination.TotalItems = len(tasksWithPaging.Items)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		tasksWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитваем количество страниц
		tasksWithPaging.Pagination.TotalPages =
			int(math.Ceil(float64(tasksWithPaging.Pagination.TotalItems) / float64(pagingOptions.PageSize)))
	}

	//Значения по умолчанию для pagingOptions
	if pagingOptions.Order != "desc" {
		pagingOptions.Order = "asc"
	}
	if pagingOptions.Page <= 0 {
		pagingOptions.Page = 1
	}
	if pagingOptions.PageSize <= 0 {
		pagingOptions.PageSize = tasksWithPaging.Pagination.TotalItems
	}
	if pagingOptions.OrderBy == "" {
		pagingOptions.OrderBy = "id"
	}
	tasksWithPaging.Pagination.Options = pagingOptions

	//Добываем выборку с учетом параметров пагинации
	err = sr.DB.Order(pagingOptions.OrderBy+" "+string(pagingOptions.Order)).
		Limit(pagingOptions.PageSize).
		Offset((pagingOptions.Page-1)*pagingOptions.PageSize).
		Find(&tasksWithPaging.Items).Where("user_id = ?, status_id = ?, category_id = ?", taskFilters.UserId, taskFilters.StatusId, taskFilters.CategoryId).Error

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
	//TODO: Обновить ненулевые поля. Если пользователь обычный юзер, то обновить может только статус СВОЕЙ задачи
	//task.Name = updatedTask.Name
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
