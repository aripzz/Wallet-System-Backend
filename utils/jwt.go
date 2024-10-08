package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(username string, user_id uint64, secret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
