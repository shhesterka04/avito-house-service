package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("your_secret_key")

const loginTime = 3 * time.Hour

func GenerateJWT(userType string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(loginTime)),
		Issuer:    "house-service",
		Subject:   userType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
