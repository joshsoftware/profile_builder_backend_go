package service

import (
	"context"

	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

type UserLoginServive interface {
	GenerateLoginToken(ctx context.Context, email string) (string, error)
}

func (userService *service) GenerateLoginToken(ctx context.Context, email string) (string, error) {
	if len((email)) == 0 {
		return "", errors.ErrEmptyPayload
	}

	userId, err := userService.UserLoginRepo.GetUserIdByEmail(ctx, email)
	if err != nil || userId == 0 {
		return "", err
	}

	token, err := userService.createToken(userId, email)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (userLoginService *service) createToken(userId int64, email string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		return "", errors.ErrSecretKey
	}

	claims := userLoginService.createClaims(userId, email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (profileSvc *service) createClaims(userId int64, email string) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userId":     userId,
		"email":      email,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
}
