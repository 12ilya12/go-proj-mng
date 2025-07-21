package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/reminder-service/gen/reminder"
	"github.com/12ilya12/go-proj-mng/repos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService struct {
	taskRepo       repos.TaskRepository
	statusRepo     repos.StatusRepository
	categoryRepo   repos.CategoryRepository
	userRepo       repos.UserRepository
	reminderClient reminder.ReminderServiceClient
}

func NewTaskService(
	taskRepo repos.TaskRepository,
	statusRepo repos.StatusRepository,
	categoryRepo repos.CategoryRepository,
	userRepo repos.UserRepository,
	reminderClient reminder.ReminderServiceClient,

) TaskService {
	return TaskService{taskRepo, statusRepo, categoryRepo, userRepo, reminderClient}
}

func (ts *TaskService) GetAll(pagingOptions pagination.PagingOptions, taskFilters common.TaskFilters) (tasksWithPaging pagination.Paging[models.Task], err error) {
	tasksWithPaging, err = ts.taskRepo.GetAll(pagingOptions, taskFilters)
	return
}

func (ts *TaskService) GetById(id uint) (task models.Task, err error) {
	task, err = ts.taskRepo.GetById(id)
	return
}

// Функция создания напоминания о задаче
func (ts *TaskService) createTaskReminder(task *models.Task, daysBefore int32) {
	_, err := ts.reminderClient.CreateReminder(context.Background(), &reminder.CreateReminderRequest{
		TaskId:     fmt.Sprintf("%d", task.ID),
		Message:    "Напоминание о скором наступлении дедлайна задачи " + task.Name,
		DaysBefore: daysBefore,
		Deadline:   timestamppb.New(task.Deadline),
	})
	if err != nil {
		log.Printf("При создании напоминания произошла ошибка: %v", err)
	}
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

	//Добавляем напоминание о задаче за 1 день до дедлайна
	ts.createTaskReminder(task, 1)

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

	updatedTask, err = ts.GetById(paramsForUpdate.ID)
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}
	//Разные возможности в зависимости от роли пользователя
	if strings.ToLower(userInfo.UserRole) == "admin" {
		//Администратор может изменять все поля задачи
		if paramsForUpdate.Name != "" {
			updatedTask.Name = paramsForUpdate.Name
		}
		if paramsForUpdate.Description != "" {
			updatedTask.Description = paramsForUpdate.Description
		}
		if paramsForUpdate.StatusId != 0 {
			updatedTask.StatusId = paramsForUpdate.StatusId
		}
		if paramsForUpdate.CategoryId != 0 {
			updatedTask.CategoryId = paramsForUpdate.CategoryId
		}
		if paramsForUpdate.UserId != 0 {
			updatedTask.UserId = paramsForUpdate.UserId
		}
		if !paramsForUpdate.Deadline.IsZero() {
			updatedTask.Deadline = paramsForUpdate.Deadline
		}
		if paramsForUpdate.Priority != 0 {
			updatedTask.Priority = paramsForUpdate.Priority
		}
	} else {
		//Пользователь не администратор, поэтому он может менять только статус СВОЕЙ задачи
		if updatedTask.UserId == uint(userInfo.UserId) {
			if paramsForUpdate.StatusId != 0 {
				updatedTask.StatusId = paramsForUpdate.StatusId
			}
		} else {
			err = common.ErrUserHasNotPermissionToEditTask
			return
		}
	}

	err = ts.taskRepo.Update(&updatedTask)
	return
}

func (ts *TaskService) Delete(id uint) (err error) {
	_, err = ts.GetById(id)
	if err != nil {
		//Не найден статус с заданным идентификатором, либо другая проблема с БД
		return
	}

	if ts.taskRepo.HasDependencies(id) {
		err = common.ErrTaskHasRelatedDependency
		return
	}
	err = ts.taskRepo.Delete(id)
	return
}
