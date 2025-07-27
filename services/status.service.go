package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type StatusService interface {
	GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status], err error)
	GetById(id uint) (status models.Status, err error)
	Create(status *models.Status) (err error)
	Update(paramsForUpdate *models.Status) (updatedStatus models.Status, err error)
	Delete(id uint) (err error)
}

type StatusServiceImpl struct {
	statusRepo repos.StatusRepository
}

func NewStatusServiceImpl(statusRepo repos.StatusRepository) StatusService {
	return &StatusServiceImpl{statusRepo}
}

func (ss *StatusServiceImpl) GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Status], err error) {
	statusesWithPaging, err = ss.statusRepo.GetAll(pagingOptions)
	return
}

func (ss *StatusServiceImpl) GetById(id uint) (status models.Status, err error) {
	status, err = ss.statusRepo.GetById(id)
	return
}

func (ss *StatusServiceImpl) Create(status *models.Status) (err error) {
	err = ss.statusRepo.Create(status)
	return
}

func (ss *StatusServiceImpl) Update(paramsForUpdate *models.Status) (updatedStatus models.Status, err error) {
	updatedStatus, err = ss.statusRepo.Update(paramsForUpdate)
	return
}

func (ss *StatusServiceImpl) Delete(id uint) (err error) {
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
