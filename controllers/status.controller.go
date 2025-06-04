package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/services"
)

type StatusController struct {
	statusService services.StatusService
}

func NewStatusController(statusService services.StatusService) StatusController {
	return StatusController{statusService}
}

func (sc *StatusController) GetAll(w http.ResponseWriter, r *http.Request) {
	//TODO: Реализовать пагинацию. Параметры пагинации будут в Query.
	var statuses []models.Status
	statuses, err := sc.statusService.GetAll( /*pagingOptions*/ )
	if err != nil {
		//TODO: Сообщение об ошибке
	}
	//TODO: В ответе помимо статусов должны быть данные пагинации
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}
