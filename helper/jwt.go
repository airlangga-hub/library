package helper

import "github.com/golang-jwt/jwt/v5"

type MyClaims struct {
	jwt.RegisteredClaims
	UserID int
}