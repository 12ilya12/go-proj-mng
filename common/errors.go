package common

import (
	"errors"
)

var ErrStatusHasRelatedTasks = errors.New("у данного статуса есть связанные задачи")
var ErrCategoryHasRelatedTasks = errors.New("у данной категории есть связанные задачи")
var ErrTaskHasRelatedDependency = errors.New("У данной задачи есть зависимости")
var ErrUserHasNotPermissionToEditTask = errors.New("Пользователь может редактировать только свою задачу")
var ErrTaskDepToItself = errors.New("Нельзя связывать задачу саму с собой")
var ErrDepOnlyBetweenUserTasks = errors.New("Пользователь может создавать зависимости только между своими задачами")
