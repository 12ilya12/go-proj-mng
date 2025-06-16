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
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) TaskController {
	return TaskController{taskService}
}

func (sc *TaskController) GetAll(w http.ResponseWriter, r *http.Request) {
	//Считываем параметры пагинации из query
	var pagingOptions pagination.PagingOptions
	utils.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	//Считываем параметры фильтрации из body
	taskFilters := common.TaskFilters{}
	err := json.NewDecoder(r.Body).Decode(&taskFilters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Дополняем параметры фильтрации информацией о пользователе
	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))
	taskFilters.UserInfo = userInfo

	tasksWithPaging, err := sc.taskService.GetAll(pagingOptions, taskFilters)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasksWithPaging)
}

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

func (sc *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: Валидация полей задачи
	err = sc.taskService.Create(&task)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (sc *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedTask := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedTask.ID = uint(id)

	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))

	err = sc.taskService.Update(&updatedTask, userInfo)
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

func (sc *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = sc.taskService.Delete(id)
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
