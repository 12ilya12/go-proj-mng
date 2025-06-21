package controllers

import (
	"encoding/json"
	"errors"
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

type StatusController struct {
	statusService services.StatusService
	vld           *validator.Validate
}

func NewStatusController(statusService services.StatusService) StatusController {
	vld := validator.New()
	return StatusController{statusService, vld}
}

// @Summary Получить все статусы
// @Description Позволяет получить список всех статусов. Доступно всем пользователям.
// @ID get-all-statuses
// @Tags Статусы
// @Produce json
// @Param Page query int false "Номер страницы"
// @Param PageSize query int false "Размер страницы"
// @Param Order query pagination.Order false "По возрастанию/по убыванию"
// @Param OrderBy query string false "Характеристика для сортировки"
// @Success 200 {object} pagination.Paging[models.Status]
// @Failure 502 {string} string
// @Router /statuses [get]
func (sc *StatusController) GetAll(w http.ResponseWriter, r *http.Request) {
	var pagingOptions pagination.PagingOptions
	utils.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	statusesWithPaging, err := sc.statusService.GetAll(pagingOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusesWithPaging)
}

// @Summary Получить статус по идентификатору
// @Description Позволяет получить статус по его идентификатору. Доступно всем пользователям.
// @ID get-status
// @Tags Статусы
// @Produce  json
// @Param id path int true "Идентификатор статуса"
// @Success 200 {object} models.Status
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Статус с заданным идентификатором не найден"
// @Failure 502 {string} string
// @Router /statuses/{id} [get]
func (sc *StatusController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status, err := sc.statusService.GetById(uint(id))
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
	json.NewEncoder(w).Encode(status)
}

// @Summary Создать статус
// @Description Позволяет создать новый статус. Доступно только для администраторов.
// @ID create-status
// @Tags Статусы
// @Produce  json
// @Success 200 {object} models.Status
// @Failure 400 {string} string "Параметры нового статуса некорректны"
// @Failure 502 {string} string
// @Router /statuses [post]
func (sc *StatusController) Create(w http.ResponseWriter, r *http.Request) {
	status := models.Status{}
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sc.vld.Struct(status)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = sc.statusService.Create(&status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// @Summary Обновить статус
// @Description Позволяет обновить данные статуса. Доступно только для администраторов.
// @ID update-status
// @Tags Статусы
// @Produce json
// @Param id path int true "Идентификатор статуса"
// @Success 200 {object} models.Status
// @Failure 400 {string} string "Параметры статуса некорректны"
// @Failure 404 {string} string "Статус с заданным идентификатором не найден"
// @Failure 502 {string} string
// @Router /statuses/{id} [patch]
func (sc *StatusController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate := models.Status{}
	err = json.NewDecoder(r.Body).Decode(&paramsForUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate.ID = uint(id)

	err = sc.vld.Struct(paramsForUpdate)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	updatedStatus, err := sc.statusService.Update(&paramsForUpdate)
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
	json.NewEncoder(w).Encode(updatedStatus)
}

// @Summary Удалить статус
// @Description Позволяет удалить статус. Доступно только для администраторов.
// @ID delete-status
// @Tags Статусы
// @Produce json
// @Param id path int true "Идентификатор статуса"
// @Success 200
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Статус с заданным идентификатором не найден"
// @Failure 409 {string} string "С удаляемым статусом есть связанные задачи"
// @Failure 502 {string} string
// @Router /statuses/{id} [delete]
func (sc *StatusController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = sc.statusService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, common.ErrStatusHasRelatedTasks) {
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
