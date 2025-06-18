package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type StatusService struct {
	statusRepo repos.StatusRepository
}

func NewStatusService(statusRepo repos.StatusRepository) StatusService {
	return StatusService{statusRepo}
}

func (ss *StatusService) GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status], err error) {
	statusesWithPaging, err = ss.statusRepo.GetAll(pagingOptions)
	return
}

func (ss *StatusService) GetById(id uint) (status models.Status, err error) {
	status, err = ss.statusRepo.GetById(id)
	return
}

func (ss *StatusService) Create(status *models.Status) (err error) {
	err = ss.statusRepo.Create(status)
	return
}

func (ss *StatusService) Update(paramsForUpdate *models.Status) (updatedStatus models.Status, err error) {
	updatedStatus, err = ss.statusRepo.Update(paramsForUpdate)
	return
}

func (ss *StatusService) Delete(id uint) (err error) {
	_, err = ss.GetById(id)
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}

	if ss.statusRepo.HasTasks(id) {
		err = common.ErrStatusHasRelatedTasks
		return
	}
	err = ss.statusRepo.Delete(id)
	return
}
