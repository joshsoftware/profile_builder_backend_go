package service

import (
	"context"

	jwttoken "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/jwt_token"
)

// UserLoginServive contains methods of creation of tokens
type UserLoginServive interface {
	GenerateLoginToken(ctx context.Context, email string) (string, error)
}

func (userService *service) GenerateLoginToken(ctx context.Context, email string) (string, error) {
	userID, err := userService.UserLoginRepo.GetUserIDByEmail(ctx, email)
	if err != nil || userID == 0 {
		return "", err
	}

	token, err := jwttoken.CreateToken(userID, email)

	if err != nil {
		return "", err
	}
	return token, nil
}
