package models

import (
	"github.com/dgrijalva/jwt-go"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}
