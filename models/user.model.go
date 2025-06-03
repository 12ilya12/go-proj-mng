package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Login    string `gorm:"type:varchar(255);not null;" json:"login"`
	Password string `gorm:"not null;" json:"password,omitempty"`
	FullName string `gorm:"type:varchar(255);not null;" json:"fullname"`
	Email    string `gorm:"type:varchar(255);not null;" json:"email"`
	Role     string `gorm:"type:varchar(255);not null;" json:"role"`
}

/*
Данные для токена JWT
*/
type Claims struct {
	UserId uint
	Role   string
	jwt.StandardClaims
}
