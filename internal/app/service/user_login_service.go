package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	jwttoken "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/jwt_token"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// UserLoginServive contains methods of creation of tokens
type UserLoginServive interface {
	GenerateLoginToken(ctx context.Context, email string) (specs.LoginResponse, error)
}

func (userService *service) GenerateLoginToken(ctx context.Context, email string) (res specs.LoginResponse, err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	userInfo, err := userService.UserLoginRepo.GetUserInfoByEmail(ctx, email)
	if err != nil || userInfo.ID == 0 || userInfo.Role == "" {
		return res, err
	}

	var profileID int
	if userInfo.Role == constants.Admin {
		profileID = 0
	} else {
		profileID, err = userService.ProfileRepo.GetProfileIdByEmail(ctx, email, tx)
		if err != nil {
			zap.S().Error("Error getting profile id by email: ", err)
			return res, err
		}
	}

	token, err := jwttoken.CreateToken(userInfo.ID, email)
	if err != nil {
		return res, err
	}

	loginResponse := specs.LoginResponse{
		ProfileID: profileID,
		Role:      userInfo.Role,
		Token:     token,
	}
	zap.S().Info("Login successful for user: ", email)
	return loginResponse, nil
}
