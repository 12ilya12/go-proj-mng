package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type TaskService struct {
	taskRepo     repos.TaskRepository
	statusRepo   repos.StatusRepository
	categoryRepo repos.CategoryRepository
	userRepo     repos.UserRepository
}

func NewTaskService(
	taskRepo repos.TaskRepository,
	statusRepo repos.StatusRepository,
	categoryRepo repos.CategoryRepository,
	userRepo repos.UserRepository,
) TaskService {
	return TaskService{taskRepo, statusRepo, categoryRepo, userRepo}
}

func (ts *TaskService) GetAll(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task], err error) {
	tasksWithPaging, err = ts.taskRepo.GetAll(pagingOptions, taskFilters)
	return
}

func (ts *TaskService) GetById(id uint) (task models.Task, err error) {
	task, err = ts.taskRepo.GetById(id)
	return
}

func (ts *TaskService) Create(task *models.Task) (err error) {
	_, err = ts.statusRepo.GetById(task.StatusId)
	if err != nil {
		//Не найден статус, заданный в задаче
		return
	}
	_, err = ts.categoryRepo.GetById(task.CategoryId)
	if err != nil {
		//Не найдена категория, заданная в задаче
		return
	}
	_, err = ts.userRepo.GetById(task.UserId)
	if err != nil {
		//Не найден пользователь, заданный в задаче
		return
	}

	err = ts.taskRepo.Create(task)
	return
}

func (ts *TaskService) Update(paramsForUpdate *models.Task, userInfo common.UserInfo) (updatedTask models.Task, err error) {
	if paramsForUpdate.StatusId != 0 {
		_, err = ts.statusRepo.GetById(paramsForUpdate.StatusId)
		if err != nil {
			//Не найден статус, заданный в задаче
			return
		}
	}
	if paramsForUpdate.CategoryId != 0 {
		_, err = ts.categoryRepo.GetById(paramsForUpdate.CategoryId)
		if err != nil {
			//Не найдена категория, заданная в задаче
			return
		}
	}
	if paramsForUpdate.UserId != 0 {
		_, err = ts.userRepo.GetById(paramsForUpdate.UserId)
		if err != nil {
			//Не найден пользователь, заданный в задаче
			return
		}
	}

	updatedTask, err = ts.taskRepo.Update(paramsForUpdate, userInfo)
	return
}

func (ts *TaskService) Delete(id uint) (err error) {
	err = ts.taskRepo.Delete(id)
	return
}
