package jwttoken

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

func CreateToken(userId int64, email string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		return "", errors.ErrSecretKey
	}

	claims := createClaims(userId, email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func createClaims(userId int64, email string) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userId":     userId,
		"email":      email,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
}
