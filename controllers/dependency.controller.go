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
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DependencyController struct {
	dependencyService services.DependencyService
}

func NewDependencyController(dependencyService services.DependencyService) DependencyController {
	return DependencyController{dependencyService}
}

func (cc *DependencyController) Get(w http.ResponseWriter, r *http.Request) {
	var pagingOptions pagination.PagingOptions
	u.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	categoriesWithPaging, err := cc.dependencyService.GetAll(pagingOptions)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoriesWithPaging)
}

/* func (cc *DependencyController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dependency, err := cc.dependencyService.GetById(id)
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
	json.NewEncoder(w).Encode(dependency)
} */

func (cc *DependencyController) Create(w http.ResponseWriter, r *http.Request) {
	dependencyDto := models.Dependency{}
	err := json.NewDecoder(r.Body).Decode(&dependencyDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: Валидация данных пользователя
	err = cc.dependencyService.Create(&dependencyDto)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dependencyDto)
}

/* func (cc *DependencyController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedDependency := models.Dependency{}
	err = json.NewDecoder(r.Body).Decode(&updatedDependency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedDependency.ID = uint(id)
	err = cc.dependencyService.Update(&updatedDependency)
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
	json.NewEncoder(w).Encode(updatedDependency)
}
*/
/* func (cc *DependencyController) HasTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hasTasks, err := cc.dependencyService.HasTasks(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	u.Respond(w, map[string]interface{}{"hasTasks": hasTasks})
} */

func (cc *DependencyController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = cc.dependencyService.Delete(id)
	if err != nil {
		if errors.Is(err, common.ErrDependencyHasRelatedTasks) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	//w.Header().Add("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(nil)
}

/* func (cc *DependencyController) DeleteForce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = cc.dependencyService.DeleteForce(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
} */
