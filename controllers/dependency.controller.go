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
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pagingOptions pagination.PagingOptions
	u.QueryDecoder.Decode(&pagingOptions, r.URL.Query())

	dependenciesWithPaging, err := cc.dependencyService.Get(taskId, pagingOptions)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dependenciesWithPaging)
}

func (cc *DependencyController) Create(w http.ResponseWriter, r *http.Request) {
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

	userInfo := common.UserInfo{}
	userInfo.UserId, _ = strconv.Atoi(fmt.Sprintf("%d", r.Context().Value(common.UserContextKey)))
	userInfo.UserRole = fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))

	//TODO: Валидация полей зависимости
	err = cc.dependencyService.Create(taskId, &dependency, userInfo)
	if err != nil {
		//TODO: Проверить какие ошибки может выдать gorm
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dependency)
}

func (cc *DependencyController) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = cc.dependencyService.Delete(taskId, dependencyId)
	if err != nil {
		/* if errors.Is(err, common.ErrDependencyHasRelatedTasks) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else  */if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//Другая проблема с БД
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
}
