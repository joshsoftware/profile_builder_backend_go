package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	jwttoken "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/jwt_token"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
)

var (
	userID0   = 0
	UserID1   = 1
	TestEmail = "test@example.com"
)

func TestUserLogin(t *testing.T) {

	mockUserLogin := new(mocks.UserStorer)
	var repodeps = service.RepoDeps{
		UserLoginDeps: mockUserLogin,
	}
	userLoginService := service.NewServices(repodeps)

	tests := []struct {
		Name          string
		Email         string
		MockSetup     func(*mocks.UserStorer, string)
		MockTokenFunc func(int64, string) (string, error)
		Expectespecsken string
		ExpectedError error
	}{
		{
			Name:  "success_of_login",
			Email: TestEmail,
			MockSetup: func(mockUserStorer *mocks.UserStorer, email string) {
				mockUserStorer.On("GetUserIdByEmail", mock.Anything, email).Return(int64(UserID1), nil).Once()
			},
			MockTokenFunc: func(id int64, email string) (string, error) {
				return "valid_token", nil
			},
			Expectespecsken: "valid_token",
			ExpectedError: nil,
		},
		{
			Name:  "failed_GetUserIdByEmail",
			Email: TestEmail,
			MockSetup: func(mockUserStorer *mocks.UserStorer, email string) {
				mockUserStorer.On("GetUserIdByEmail", mock.Anything, email).Return(int64(userID0), errors.New("repository error")).Once()
			},
			MockTokenFunc: nil,
			Expectespecsken: "",
			ExpectedError: errors.New("repository error"),
		},
		{
			Name:  "CreateToken_error",
			Email: TestEmail,
			MockSetup: func(mockUserStorer *mocks.UserStorer, email string) {
				mockUserStorer.On("GetUserIdByEmail", mock.Anything, email).Return(int64(UserID1), nil).Once()
			},
			MockTokenFunc: func(id int64, email string) (string, error) {
				return "", errors.New("token creation error")
			},
			Expectespecsken: "",
			ExpectedError: errors.New("token creation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockSetup(mockUserLogin, tt.Email)

			if tt.MockTokenFunc != nil {
				patch, _ := mpatch.PatchMethod(jwttoken.CreateToken, tt.MockTokenFunc)
				defer patch.Unpatch()
			}

			token, err := userLoginService.GenerateLoginToken(context.Background(), tt.Email)
			assert.Equal(t, tt.Expectespecsken, token)
			assert.Equal(t, tt.ExpectedError, err)

			mockUserLogin.AssertExpectations(t)
		})

	}
}
