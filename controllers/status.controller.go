package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/services"
	"github.com/gorilla/mux"
)

type StatusController struct {
	statusService services.StatusService
}

func NewStatusController(statusService services.StatusService) StatusController {
	return StatusController{statusService}
}

func (sc *StatusController) GetAll(w http.ResponseWriter, r *http.Request) {
	//TODO: Реализовать пагинацию. Параметры пагинации будут в Query.
	statuses, err := sc.statusService.GetAll( /*pagingOptions*/ )
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: В ответе помимо статусов должны быть данные пагинации
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (sc *StatusController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status, err := sc.statusService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (sc *StatusController) Create(w http.ResponseWriter, r *http.Request) {
	statusDto := models.Status{}
	err := json.NewDecoder(r.Body).Decode(&statusDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: Валидация данных пользователя
	err = sc.statusService.Create(&statusDto)
	if err != nil {
		//TODO: Статус ответа должен определяться в зависимости от ошибки.
		//Не факт, что проблема в запросе. Например, могут быть проблемы с БД.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusDto)
}

func (sc *StatusController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newStatusDto := models.Status{}
	err = json.NewDecoder(r.Body).Decode(&newStatusDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = sc.statusService.Update(id, &newStatusDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newStatusDto)
}

func (sc *StatusController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = sc.statusService.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}
