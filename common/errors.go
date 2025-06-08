package common

import (
	"errors"
)

var ErrStatusHasRelatedTasks = errors.New("У данного статуса есть связанные задачи")
