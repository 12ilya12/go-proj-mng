package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	//ID       uint32 `gorm:"primary_key" json:"id"`
	gorm.Model
	Login    string `gorm:"type:varchar(255);not null;" json:"login"`
	Password string `gorm:"not null;" json:"password"`
	FullName string `gorm:"type:varchar(255);not null;" json:"fullname"`
	Email    string `gorm:"type:varchar(255);not null;" json:"email"`
	Role     string `gorm:"type:varchar(255);not null;" json:"role"`
}

/*type UserResponse struct {
	ID       uint32 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserCreate struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}*/

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}
