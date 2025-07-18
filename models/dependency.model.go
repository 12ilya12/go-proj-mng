package models

type Dependency struct {
	ID           uint `gorm:"primary_key" json:"id"`
	ParentTask   Task `gorm:"foreignkey:ParentTaskId;association_foreignkey:ID" json:"-"`
	ParentTaskId uint `gorm:"not null" json:"parent_task_id" validate:"required"`
	ChildTask    Task `gorm:"foreignkey:ChildTaskId;association_foreignkey:ID" json:"-"`
	ChildTaskId  uint `gorm:"not null" json:"child_task_id" validate:"required"`
}
