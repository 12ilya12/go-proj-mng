package services

import (
	"strings"

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

func (ds *DependencyService) Get(parentTaskId uint, pagingOptions pagination.PagingOptions) (dependenciesWithPaging pagination.Paging[models.Dependency], err error) {
	dependenciesWithPaging, err = ds.dependencyRepo.Get(parentTaskId, pagingOptions)
	return
}

func (ds *DependencyService) Create(dependency *models.Dependency, userInfo common.UserInfo) (err error) {
	parentTask, err := ds.taskRepo.GetById(dependency.ParentTaskId)
	if err != nil {
		//Родительская задача не найдена
		return
	}
	childTask, err := ds.taskRepo.GetById(dependency.ChildTaskId)
	if err != nil {
		//Дочерняя задача не найдена
		return
	}
	if dependency.ParentTaskId == dependency.ChildTaskId {
		err = common.ErrTaskDepToItself
		return
	}
	//Обычный пользователь может создавать зависимости только между своим задачами
	if strings.ToLower(userInfo.UserRole) == "user" &&
		(userInfo.UserId != int(parentTask.UserId) || userInfo.UserId != int(childTask.UserId)) {
		err = common.ErrDepOnlyBetweenUserTasks
		return
	}

	err = ds.dependencyRepo.Create(dependency, userInfo, parentTask.UserId, childTask.UserId)
	return
}

func (ds *DependencyService) Delete(parentTaskId uint, dependencyId uint) (err error) {
	err = ds.dependencyRepo.Delete(parentTaskId, dependencyId)
	return
}
