package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type TaskService struct {
	taskRepo repos.TaskRepository
}

func NewTaskService(taskRepo repos.TaskRepository) TaskService {
	return TaskService{taskRepo}
}

func (ss *TaskService) GetAll(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task], err error) {
	if taskFilters.UserInfo.UserRole == "ADMIN" {
		tasksWithPaging, err = ss.taskRepo.GetAll(pagingOptions, taskFilters)
	} else {
		tasksWithPaging, err = ss.taskRepo.GetAllForUser(pagingOptions, taskFilters)
	}

	return
}

func (ss *TaskService) GetById(id int) (task models.Task, err error) {
	task, err = ss.taskRepo.GetById(id)
	return
}

func (ss *TaskService) Create(task *models.Task) (err error) {
	err = ss.taskRepo.Create(task)
	return
}

func (ss *TaskService) Update(task *models.Task, userInfo common.UserInfo) (err error) {
	err = ss.taskRepo.Update(task, userInfo)
	return
}

func (ss *TaskService) Delete(id int) (err error) {
	err = ss.taskRepo.Delete(id)
	return
}
