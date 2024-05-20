package profile

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type service struct {
	Repo repository.ProfileStorer
}

type Service interface {
	CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error
	GenerateLoginToken(ctx context.Context, email string) (string, error)
}

func NewServices(repo repository.ProfileStorer) Service {
	return &service{
		Repo: repo,
	}
}

func (profileSvc *service) CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error {
	err := profileSvc.Repo.CreateProfile(profileDetail, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (profileSvc *service) GenerateLoginToken(ctx context.Context, email string) (string, error) {
	userId, err := profileSvc.Repo.GetProfileByEmail(ctx, email)
	if err != nil || userId == 0 {
		return "", err
	}
	token, err := profileSvc.createToken(email, userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (profileSvc *service) createClaims(email string, userId int64) jwt.MapClaims {
	return jwt.MapClaims{
		"authorised": true,
		"userId":     userId,
		"email":      email,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
}

func (profileSvc *service) createToken(email string, userId int64) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("secret key not found")
	}
	claims := profileSvc.createClaims(email, userId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
