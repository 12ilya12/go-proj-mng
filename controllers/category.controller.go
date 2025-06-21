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
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CategoryController struct {
	categoryService services.CategoryService
	vld             *validator.Validate
}

func NewCategoryController(categoryService services.CategoryService) CategoryController {
	vld := validator.New()
	return CategoryController{categoryService, vld}
}

// @Summary Получить все категории
// @Description Позволяет получить список всех категорий. Доступно всем пользователям.
// @ID get-all-categories
// @Tags Категории
// @Produce json
// @Param Page query int false "Номер страницы"
// @Param PageSize query int false "Размер страницы"
// @Param Order query pagination.Order false "По возрастанию/по убыванию"
// @Param OrderBy query string false "Характеристика для сортировки"
// @Success 200 {object} pagination.Paging[models.Category]
// @Failure 502 {string} string
// @Router /categories [get]
func (cc *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	var pagingOptions pagination.PagingOptions
	u.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	categoriesWithPaging, err := cc.categoryService.GetAll(pagingOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoriesWithPaging)
}

// @Summary Получить категорию по идентификатору
// @Description Позволяет получить категорию по его идентификатору. Доступно всем пользователям.
// @ID get-category
// @Tags Категории
// @Produce  json
// @Param id path int true "Идентификатор категории"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Категория с заданным идентификатором не найден"
// @Failure 502 {string} string
// @Router /categories/{id} [get]
func (cc *CategoryController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := cc.categoryService.GetById(uint(id))
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
	json.NewEncoder(w).Encode(category)
}

// @Summary Создать категорию
// @Description Позволяет создать новую категорию. Доступно только для администраторов.
// @ID create-category
// @Tags Категории
// @Produce  json
// @Success 200 {object} models.Category
// @Failure 400 {string} string "Параметры новой категории некорректны"
// @Failure 502 {string} string
// @Router /categories [post]
func (cc *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cc.vld.Struct(category)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = cc.categoryService.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// @Summary Обновить категорию
// @Description Позволяет обновить данные категории. Доступно только для администраторов.
// @ID update-category
// @Tags Категории
// @Produce json
// @Param id path int true "Идентификатор категории"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "Параметры категории некорректны"
// @Failure 404 {string} string "Категория с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /categories/{id} [patch]
func (cc *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate := models.Category{}
	err = json.NewDecoder(r.Body).Decode(&paramsForUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	paramsForUpdate.ID = uint(id)

	err = cc.vld.Struct(paramsForUpdate)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	updatedCategory, err := cc.categoryService.Update(&paramsForUpdate)
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
	json.NewEncoder(w).Encode(updatedCategory)
}

type HasTasks struct {
	HasTasks bool `json:"hasTasks"`
}

// @Summary Есть ли у категории связанные задачи
// @Description Позволяет получить информацию о том, есть ли связанные задачи с категорией по ID. Доступно для всех пользователей.
// @ID has-tasks-category
// @Tags Категории
// @Produce json
// @Param id path int true "Идентификатор категории"
// @Success 200 {object} controllers.HasTasks
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Категория с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /categories/{id}/has-tasks [get]
func (cc *CategoryController) HasTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hasTasks, err := cc.categoryService.HasTasks(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	hasTasksResp := HasTasks{HasTasks: hasTasks}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasTasksResp)
}

// @Summary Удалить категорию
// @Description Позволяет удалить категорию. Доступно только для администраторов.
// @ID delete-category
// @Tags Категории
// @Produce json
// @Param id path int true "Идентификатор категории"
// @Success 200
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Категория с заданным идентификатором не найдена"
// @Failure 409 {string} string "С удаляемой категорией есть связанные задачи"
// @Failure 502 {string} string
// @Router /categories/{id} [delete]
func (cc *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = cc.categoryService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, common.ErrCategoryHasRelatedTasks) {
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

// @Summary Принудительно удалить категорию со связанными задачами
// @Description Позволяет принудительно удалить категорию по идентификатору вместе со всеми связаннами задачами каскадно. Доступно только для администраторов.
// @ID delete-force-category
// @Tags Категории
// @Produce json
// @Param id path int true "Идентификатор категории"
// @Success 200
// @Failure 400 {string} string "Некорректный идентификатор"
// @Failure 404 {string} string "Категория с заданным идентификатором не найдена"
// @Failure 502 {string} string
// @Router /categories/{id}/force [delete]
func (cc *CategoryController) DeleteForce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = cc.categoryService.DeleteForce(uint(id))
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
