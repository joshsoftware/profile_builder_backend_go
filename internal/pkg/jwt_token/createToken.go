package jwttoken

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// CreateToken used to generate a token
func CreateToken(userID int64, email string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		return "", errors.ErrSecretKey
	}

	expirationHoursStr := os.Getenv("TOKEN_EXPIRATION_HOURS")
	if expirationHoursStr == "" {
		zap.S().Fatal("TOKEN_EXPIRATION_HOURS is not set")
		return "", errors.ErrTokenExpirationHours
	}
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		log.Fatalf("Error parsing TOKEN_EXPIRATION_HOURS: %v", err)
		return "", err
	}

	claims := createClaims(userID, email, time.Duration(expirationHours)*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Error in generating token: %v", err)
		return "", err
	}
	return tokenString, nil
}

// CreateClaims to generate claims that are required to create token
func createClaims(userID int64, email string, expiration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userID":     userID,
		"email":      email,
		"exp":        time.Now().Add(expiration).Unix(),
	}
}
