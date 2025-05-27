package models

import (
	"github.com/google/uuid"
)

type Dependency struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ParentTask   Task      `gorm:"foreignkey:ParentTaskId;association_foreignkey:ID"`
	ParentTaskId uint32    `gorm:"not null" json:"parent_task_id"`
	ChildTask    Task      `gorm:"foreignkey:ChildTaskId;association_foreignkey:ID"`
	ChildTaskId  uint32    `gorm:"not null" json:"child_task_id"`
}
