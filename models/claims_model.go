package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}
