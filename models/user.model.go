package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Login    string    `gorm:"type:varchar(255);not null;"`
	Password string    `gorm:"not null;"`
	FullName string    `gorm:"type:varchar(255);not null;"`
	Email    string    `gorm:"type:varchar(255);not null;"`
	Role     string    `gorm:"type:varchar(255);not null;"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	FullName string    `json:"fullname"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}
