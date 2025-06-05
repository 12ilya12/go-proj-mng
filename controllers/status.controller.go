package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	//var statuses []models.Status
	statuses, err := sc.statusService.GetAll( /*pagingOptions*/ )
	if err != nil {
		//TODO: Сообщение об ошибке
	}
	//TODO: В ответе помимо статусов должны быть данные пагинации
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (sc *StatusController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {

	}
	status, err := sc.statusService.GetById(id)
	if err != nil {
		//TODO: Сообщение об ошибке
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
