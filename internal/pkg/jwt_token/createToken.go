package jwttoken

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

func CreateToken(userId int64, email string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		return "", errors.ErrSecretKey
	}

	expirationHoursStr := os.Getenv("TOKEN_EXPIRATION_HOURS")
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		log.Fatalf("Error parsing TOKEN_EXPIRATION_HOURS: %v", err)
		return "", err
	}

	claims := createClaims(userId, email, time.Duration(expirationHours)*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Error in generating token: %v", err)
		return "", err
	}
	return tokenString, nil
}

func createClaims(userId int64, email string, expiration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userId":     userId,
		"email":      email,
		"exp":        time.Now().Add(expiration).Unix(),
	}
}
