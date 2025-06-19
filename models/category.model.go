package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Category struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(255);not null;" json:"name" validate:""`
}

// Хук для каскадного удаления задач, связанных с категорией
func (c *Category) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("category_id = ?", c.ID).Delete(&Task{})
	return
}
