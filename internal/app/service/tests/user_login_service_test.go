package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	jwttoken "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/jwt_token"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
)

var (
	UserID0           = 0
	UserID1           = 1
	TestAdminEmail    = "admin@example.com"
	TestEmployeeEmail = "employee@example.com"
)

func TestUserLogin(t *testing.T) {
	mockUserLogin := new(mocks.UserStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:   mockProfileRepo,
		UserLoginDeps: mockUserLogin,
	}
	userLoginService := service.NewServices(repodeps)

	mockAdminInfo := repository.User{
		ID:    1,
		Email: TestAdminEmail,
		Role:  constants.Admin,
	}

	mockEmployeeInfo := repository.User{
		ID:    2,
		Email: TestEmployeeEmail,
		Role:  constants.Employee,
	}

	mockAdminResponse := specs.LoginResponse{
		ProfileID: constants.AdminProfileID,
		Role:      constants.Admin,
		Token:     "valid_admin_token",
	}

	mockEmployeeResponse := specs.LoginResponse{
		ProfileID: 2,
		Role:      constants.Employee,
		Token:     "valid_employee_token",
	}

	tests := []struct {
		Name             string
		Email            string
		Role             string
		MockSetup        func(*mocks.UserStorer, *mocks.ProfileStorer, string, string)
		MockTokenFunc    func(int64, int, string, string) (string, error)
		ExpectedResponse specs.LoginResponse
		ExpectedError    error
	}{
		{
			Name:  "success_login_admin",
			Email: TestAdminEmail,
			Role:  constants.Admin,
			MockSetup: func(mockUserStorer *mocks.UserStorer, profileMock *mocks.ProfileStorer, email string, role string) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockUserStorer.On("GetUserInfo", mock.Anything, specs.UserInfoFilter{Email: email}).Return(mockAdminInfo, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(nil).Once()
			},
			MockTokenFunc: func(userID int64, profileID int, email string, role string) (string, error) {
				return "valid_admin_token", nil
			},
			ExpectedResponse: mockAdminResponse,
			ExpectedError:    nil,
		},
		{
			Name:  "success_employee_login",
			Email: TestEmployeeEmail,
			Role:  constants.Employee,
			MockSetup: func(mockUserStorer *mocks.UserStorer, profileMock *mocks.ProfileStorer, email, role string) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockUserStorer.On("GetUserInfo", mock.Anything, specs.UserInfoFilter{Email: email}).Return(mockEmployeeInfo, nil).Once()
				profileMock.On("GetProfileIDByEmail", mock.Anything, email, mock.Anything).Return(2, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(nil).Once()
			},
			MockTokenFunc: func(userID int64, profileID int, role, email string) (string, error) {
				return "valid_employee_token", nil
			},
			ExpectedResponse: mockEmployeeResponse,
			ExpectedError:    nil,
		},
		{
			Name:  "failed_get_user_info",
			Email: TestAdminEmail,
			MockSetup: func(mockUserStorer *mocks.UserStorer, profileMock *mocks.ProfileStorer, email, role string) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockUserStorer.On("GetUserInfo", mock.Anything, specs.UserInfoFilter{Email: email}).Return(repository.User{}, errors.New("repository error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, errors.New("repository error")).Return(errors.New("repository error")).Once()
			},
			MockTokenFunc:    nil,
			ExpectedResponse: specs.LoginResponse{},
			ExpectedError:    errors.New("repository error"),
		},
		{
			Name:  "failed_get_profile_id",
			Email: TestEmployeeEmail,
			Role:  constants.Employee,
			MockSetup: func(mockUserStorer *mocks.UserStorer, profileMock *mocks.ProfileStorer, email, role string) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockUserStorer.On("GetUserInfo", mock.Anything, specs.UserInfoFilter{Email: email}).Return(mockEmployeeInfo, nil).Once()
				profileMock.On("GetProfileIDByEmail", mock.Anything, email, mock.Anything).Return(0, errors.New("profile id error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, errors.New("profile id error")).Return(errors.New("profile id error")).Once()
			},
			MockTokenFunc:    nil,
			ExpectedResponse: specs.LoginResponse{},
			ExpectedError:    errors.New("profile id error"),
		},
		{
			Name:  "CreateToken_error",
			Email: TestAdminEmail,
			MockSetup: func(mockUserStorer *mocks.UserStorer, profileMock *mocks.ProfileStorer, email, role string) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockUserStorer.On("GetUserInfo", mock.Anything, specs.UserInfoFilter{Email: email}).Return(mockAdminInfo, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, errors.New("token creation error")).Return(errors.New("token creation error")).Once()
			},
			MockTokenFunc: func(userID int64, profileID int, email, role string) (string, error) {
				return "", errors.New("token creation error")
			},
			ExpectedResponse: specs.LoginResponse{},
			ExpectedError:    errors.New("token creation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockSetup(mockUserLogin, mockProfileRepo, tt.Email, tt.Role)

			if tt.MockTokenFunc != nil {
				patch, _ := mpatch.PatchMethod(jwttoken.CreateToken, tt.MockTokenFunc)
				defer patch.Unpatch()
			}

			token, err := userLoginService.GenerateLoginToken(context.Background(), specs.UserInfoFilter{Email: tt.Email})
			assert.Equal(t, tt.ExpectedResponse, token)
			assert.Equal(t, tt.ExpectedError, err)

			mockUserLogin.AssertExpectations(t)
			mockProfileRepo.AssertExpectations(t)
		})

	}
}

func TestRemoveToken(t *testing.T) {
	mockUserLogin := new(mocks.UserStorer)
	var repodeps = service.RepoDeps{
		UserLoginDeps: mockUserLogin,
	}
	userLoginService := service.NewServices(repodeps)
	tests := []struct {
		name        string
		token       string
		setupTokens map[string]struct{}
		expectedErr error
	}{
		{
			name:        "Token exists",
			token:       "validToken",
			setupTokens: map[string]struct{}{"validToken": {}},
			expectedErr: nil,
		},
		{
			name:        "Token does not exist",
			token:       "invalidToken",
			setupTokens: map[string]struct{}{},
			expectedErr: errs.ErrTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helpers.TokenList = tt.setupTokens
			err := userLoginService.RemoveToken(tt.token)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil {
				if _, found := helpers.TokenList[tt.token]; found {
					t.Errorf("expected token %v to be removed", tt.token)
				}
			}
		})
	}
}
