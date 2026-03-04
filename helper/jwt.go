package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID int
}

func MakeJWT(userID int, email string, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   email,
		},
		UserID: userID,
	})
	
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("helper.MakeJWT: %w", err)
	}
	
	return tokenStr, nil
}
