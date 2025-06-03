package models

type Dependency struct {
	ID           uint `gorm:"primary_key" json:"id"`
	ParentTask   Task `gorm:"foreignkey:ParentTaskId;association_foreignkey:ID"`
	ParentTaskId uint `gorm:"not null" json:"parent_task_id"`
	ChildTask    Task `gorm:"foreignkey:ChildTaskId;association_foreignkey:ID"`
	ChildTaskId  uint `gorm:"not null" json:"child_task_id"`
}
