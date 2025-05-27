package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Status      Status    `gorm:"foreignkey:StatusId;association_foreignkey:ID"`
	StatusId    uint32    `gorm:"not null" json:"status_id"`
	Category    Category  `gorm:"foreignkey:CategoryId;association_foreignkey:ID;constraint:OnDelete:CASCADE"`
	CategoryId  uint32    `gorm:"not null" json:"category_id"`
	User        User      `gorm:"foreignkey:UserId;association_foreignkey:ID"`
	UserId      uint32    `json:"user_id"`
	Deadline    time.Time `json:"deadline"`
	Priority    int       `gorm:"not null;default:1" json:"priority"`
}
