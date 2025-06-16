package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type DependencyService struct {
	dependencyRepo repos.DependencyRepository
	taskRepo       repos.TaskRepository
}

func NewDependencyService(dependencyRepo repos.DependencyRepository, taskRepo repos.TaskRepository) DependencyService {
	return DependencyService{dependencyRepo, taskRepo}
}

func (ss *DependencyService) Get(parentTaskId int, pagingOptions pagination.PagingOptions) (dependenciesWithPaging pagination.Paging[models.Dependency], err error) {
	dependenciesWithPaging, err = ss.dependencyRepo.Get(parentTaskId, pagingOptions)
	return
}

func (ss *DependencyService) Create(parentTaskId int, dependency *models.Dependency, userInfo common.UserInfo) (err error) {
	err = ss.dependencyRepo.Create(parentTaskId, dependency, userInfo)
	return
}

func (ss *DependencyService) Delete(parentTaskId int, dependencyId int) (err error) {
	err = ss.dependencyRepo.Delete(parentTaskId, dependencyId)
	return
}
