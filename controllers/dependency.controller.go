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
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DependencyController struct {
	dependencyService services.DependencyService
	vld               *validator.Validate
}

func NewDependencyController(dependencyService services.DependencyService) DependencyController {
	vld := validator.New()
	return DependencyController{dependencyService, vld}
}

// @Summary Получить все зависимости задачи
// @Description Позволяет получить список всех зависимостей задачи. Доступно всем пользователям.
// @ID get-dependencies
// @Tags Зависимости задач
// @Produce json
// @Param taskId path int true "Идентификатор задачи"
// @Param Page query int false "Номер страницы"
// @Param PageSize query int false "Размер страницы"
// @Param Order query pagination.Order false "По возрастанию/по убыванию"
// @Param OrderBy query string false "Характеристика для сортировки"
// @Success 200 {object} pagination.Paging[models.Dependency]
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Задача с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /tasks/{taskId}/dependencies [get]
func (dc *DependencyController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pagingOptions pagination.PagingOptions
	u.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	dependenciesWithPaging, err := dc.dependencyService.Get(uint(taskId), pagingOptions)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dependenciesWithPaging)
}

// @Summary Добавить зависимость у задачи
// @Description Позволяет добавить зависимость у задачи. Доступно всем пользователям, но пользователи могут создавать зависимости только между своими задачами.
// @ID create-dependency
// @Tags Зависимости задач
// @Produce json
// @Param taskId path int true "Идентификатор задачи"
// @Success 200 {object} models.Dependency
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Задача с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /tasks/{taskId}/dependencies [post]
func (dc *DependencyController) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dependency := models.Dependency{}
	err = json.NewDecoder(r.Body).Decode(&dependency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dependency.ParentTaskId = uint(taskId)

	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))

	err = dc.vld.Struct(dependency)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = dc.dependencyService.Create(&dependency, userInfo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if errors.Is(err, common.ErrTaskDepToItself) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, common.ErrDepOnlyBetweenUserTasks) {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dependency)
}

// @Summary Удалить зависимость
// @Description Позволяет удалить зависимость у задачи по идентификатору. Доступно всем пользователям.
// @ID delete-dependency
// @Tags Зависимости задач
// @Produce json
// @Param taskId path int true "Идентификатор задачи"
// @Param dependencyId path int true "Идентификатор зависимости"
// @Success 200
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Задача или зависимость с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /tasks/{taskId}/dependencies/{dependencyId} [delete]
func (dc *DependencyController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dependencyId, err := strconv.Atoi(vars["dependencyId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dc.dependencyService.Delete(uint(taskId), uint(dependencyId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
}
