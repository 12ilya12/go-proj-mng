package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/services"
	"github.com/12ilya12/go-proj-mng/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TaskController struct {
	taskService services.TaskService
	vld         *validator.Validate
}

func NewTaskController(taskService services.TaskService) TaskController {
	vld := validator.New()
	return TaskController{taskService, vld}
}

// @Summary Получить все задачи
// @Description Позволяет получить список всех задач. Доступно всем пользователям. Администратор получает весь список задач. Обычный пользователь только свои задачи.
// @ID get-all-tasks
// @Tags Задачи
// @Produce json
// @Param Page query int false "Номер страницы"
// @Param PageSize query int false "Размер страницы"
// @Param Order query pagination.Order false "По возрастанию/по убыванию"
// @Param OrderBy query string false "Характеристика для сортировки"
// @Success 200 {object} pagination.Paging[models.Task]
// @Failure 502 {string} string
// @Router /tasks [get]
func (sc *TaskController) GetAll(w http.ResponseWriter, r *http.Request) {
	//Считываем параметры пагинации из query
	var pagingOptions pagination.PagingOptions
	utils.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	//Считываем необязательные параметры фильтрации из body
	taskFilters := common.TaskFilters{}
	json.NewDecoder(r.Body).Decode(&taskFilters)

	//Дополняем параметры фильтрации информацией о пользователе
	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))
	taskFilters.UserInfo = userInfo

	tasksWithPaging, err := sc.taskService.GetAll(pagingOptions, taskFilters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksWithPaging)
}

// @Summary Получить задачу по идентификатору
// @Description Позволяет получить задачу по его идентификатору. Доступно всем пользователям.
// @ID get-task
// @Tags Задачи
// @Produce  json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Задача с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /tasks/{id} [get]
func (sc *TaskController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := sc.taskService.GetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// @Summary Создать задачу
// @Description Позволяет создать новую задачу. Доступно только для администраторов.
// @ID create-task
// @Tags Задачи
// @Produce  json
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Параметры новой задачи некорректны"
// @Failure 502 {string} string
// @Router /tasks [post]
func (sc *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sc.vld.Struct(task)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = sc.taskService.Create(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// @Summary Обновить задачу
// @Description Позволяет обновить данные задачи. Доступно для администраторов и пользователей, но пользователь может изменить только статус своей задачи.
// @ID update-task
// @Tags Задачи
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Параметры задачи некорректны"
// @Failure 404 {string} string "Задача с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /tasks/{id} [patch]
func (sc *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&paramsForUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate.ID = uint(id)

	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))

	updatedTask, err := sc.taskService.Update(&paramsForUpdate, userInfo)
	if err != nil {
		if errors.Is(err, common.ErrUserHasNotPermissionToEditTask) {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// @Summary Удалить задачу
// @Description Позволяет удалить задачу. Доступно только для администраторов.
// @ID delete-task
// @Tags Задачи
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Задача с заданным идентификатором не найдена"
// @Failure 409 {string} string "Удаляемая задача связана с другой задачей"
// @Failure 502 {string} string
// @Router /tasks/{id} [delete]
func (sc *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = sc.taskService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, common.ErrTaskHasRelatedDependency) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
}
