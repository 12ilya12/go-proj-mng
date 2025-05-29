package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       uint32 `gorm:"primary_key"`
	Login    string `gorm:"type:varchar(255);not null;"`
	Password string `gorm:"not null;"`
	FullName string `gorm:"type:varchar(255);not null;"`
	Email    string `gorm:"type:varchar(255);not null;"`
	Role     string `gorm:"type:varchar(255);not null;"`
}

type UserResponse struct {
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
}

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}
