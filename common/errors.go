package common

import (
	"errors"
)

var ErrStatusHasRelatedTasks = errors.New("у данного статуса есть связанные задачи")
var ErrCategoryHasRelatedTasks = errors.New("у данной категории есть связанные задачи")
var ErrTaskHasRelatedDependency = errors.New("У данной задачи есть зависимости")
