package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name" validator:"required"`
	Description string    `gorm:"not null" json:"description" validator:"required"`
	Status      Status    `gorm:"foreignkey:StatusId;association_foreignkey:ID" json:"-"`
	StatusId    uint      `gorm:"not null" json:"status_id" validator:"required"`
	Category    Category  `gorm:"foreignkey:CategoryId;association_foreignkey:ID;constraint:OnDelete:CASCADE" json:"-"`
	CategoryId  uint      `gorm:"not null" json:"category_id" validator:"required"`
	User        User      `gorm:"foreignkey:UserId;association_foreignkey:ID" json:"-"`
	UserId      uint      `json:"user_id" validator:"required"`
	Deadline    time.Time `json:"deadline"`
	Priority    int       `gorm:"not null;default:1" json:"priority" validator:"gte=1,lte=4"`
}
