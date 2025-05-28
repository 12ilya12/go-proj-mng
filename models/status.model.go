package models

type Status struct {
	ID   uint32 `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(255);not null;" json:"name"`
}
