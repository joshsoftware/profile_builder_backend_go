package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	jwttoken "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/jwt_token"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// UserLoginServive contains methods of creation of tokens
type UserLoginServive interface {
	GenerateLoginToken(ctx context.Context, filter specs.UserInfoFilter) (specs.LoginResponse, error)
	RemoveToken(token string) error
}

func (userService *service) GenerateLoginToken(ctx context.Context, filter specs.UserInfoFilter) (res specs.LoginResponse, err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	userInfo, err := userService.UserLoginRepo.GetUserInfo(ctx, filter)
	if err != nil || userInfo.ID == 0 || userInfo.Role == "" {
		return specs.LoginResponse{}, err
	}

	var profileID int
	if userInfo.Role == constants.Admin {
		profileID = constants.AdminProfileID
	} else {
		profileID, err = userService.ProfileRepo.GetProfileIDByEmail(ctx, filter.Email, tx)
		if err != nil {
			zap.S().Errorf("Error getting profile id : %v by email : %s ", err, filter.Email)
			return specs.LoginResponse{}, err
		}
	}

	token, err := jwttoken.CreateToken(userInfo.ID, profileID, userInfo.Role, filter.Email)
	if err != nil {
		return specs.LoginResponse{}, err
	}

	helpers.WhiteListMutext.Lock()
	helpers.TokenList[token] = struct{}{}
	helpers.WhiteListMutext.Unlock()

	loginResponse := specs.LoginResponse{
		ProfileID: profileID,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Role:      userInfo.Role,
		Token:     token,
	}
	zap.S().Infof("Login successful for user: %s", filter.Email)
	return loginResponse, nil
}

func (userService *service) RemoveToken(token string) error {
	helpers.WhiteListMutext.Lock()
	defer helpers.WhiteListMutext.Unlock()
	if _, found := helpers.TokenList[token]; !found {
		zap.S().Error("Token not found in whitelist")
		return errors.ErrTokenNotFound
	}
	delete(helpers.TokenList, token)

	zap.S().Info("Logout successfully")
	return nil
}
