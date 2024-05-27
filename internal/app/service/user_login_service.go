package service

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type userLoginService struct {
	UserLoginRepo repository.UserStorer
}

type UserLoginService interface {
	GenerateLoginToken(ctx context.Context, email string) (string, error)
}

func NewUserLoginService(db *pgx.Conn) UserLoginService {
	userLoginRepo := repository.NewUserLoginRepo(db)

	return &userLoginService{
		UserLoginRepo: userLoginRepo,
	}
}

func (userService *userLoginService) GenerateLoginToken(ctx context.Context, email string) (string, error) {
	if len((email)) == 0 {
		return "", errors.ErrEmptyPayload
	}

	userId, err := userService.UserLoginRepo.GetProfileByEmail(ctx, email)
	if err != nil || userId == 0 {
		return "", err
	}

	token, err := userService.createToken(userId, email)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (userLoginService *userLoginService) createToken(userId int64, email string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		return "", errors.ErrSectetKeyNotFound
	}

	claims := userLoginService.createClaims(userId, email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (profileSvc *userLoginService) createClaims(userId int64, email string) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userId":     userId,
		"email":      email,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
}
