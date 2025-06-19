package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Login    string `gorm:"type:varchar(255);not null;" json:"login" validator:"required"`
	Password string `gorm:"not null;" json:"password,omitempty" validator:"required"`
	FullName string `gorm:"type:varchar(255);not null;" json:"fullname"`
	Email    string `gorm:"type:varchar(255);not null;" json:"email" validator:"email"`
	Role     string `gorm:"type:varchar(255);not null;" json:"role" validator:"oneof=admin user"`
}

/*
Данные для токена JWT
*/
type Claims struct {
	UserId uint
	Role   string
	jwt.StandardClaims
}
