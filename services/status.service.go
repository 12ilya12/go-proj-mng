package services

import (
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/repos"
)

type StatusService struct {
	statusRepo repos.StatusRepository
}

func NewStatusService(statusRepo repos.StatusRepository) StatusService {
	return StatusService{statusRepo}
}

func (ss *StatusService) GetAll( /*pagingOptions*/ ) (statuses []models.Status, err error) {
	statuses, err = ss.statusRepo.GetAll( /*pagingOptions*/ )
	return
}

func (ss *StatusService) GetById(id int) (status models.Status, err error) {
	status, err = ss.statusRepo.GetById(id)
	return
}
