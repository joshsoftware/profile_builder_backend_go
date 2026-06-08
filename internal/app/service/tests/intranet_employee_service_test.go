package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	intranetMocks "github.com/joshsoftware/profile_builder_backend_go/internal/client/intranet/mocks"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	repoMocks "github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetIntranetEmployee(t *testing.T) {
	mockIntranetClient := new(intranetMocks.IntranetClient)
	mockProfileRepo := new(repoMocks.ProfileStorer)
	repodeps := service.RepoDeps{
		IntranetClient: mockIntranetClient,
		ProfileDeps:    mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		employeeID      string
		setup           func(clientMock *intranetMocks.IntranetClient, repoMock *repoMocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    specs.IntranetEmployeeResponse
		expectedError   error
	}{
		{
			name:       "Success_get_intranet_employee",
			employeeID: "EMP123",
			setup: func(clientMock *intranetMocks.IntranetClient, repoMock *repoMocks.ProfileStorer) {
				repoMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				repoMock.On("GetProfileIDByEmployeeID", mock.Anything, "EMP123", mock.Anything).Return(0, errs.ErrNoRecordFound).Once()
				repoMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				clientMock.On("GetEmployeeByID", mock.Anything, "EMP123").Return(&specs.IntranetEmployee{
					EmployeeID:        "EMP123",
					Email:             "emp@example.com",
					Name:              "John Doe",
					MobileNumber:      "1234567890",
					Gender:            "Male",
					YearsOfExperience: 5.5,
					Designation:       "Software Engineer",
					JoshDOJ:           "2020-01-01",
					LinkedinURL:       "https://linkedin.com/in/johndoe",
					GithubURL:         "https://github.com/johndoe",
					PrimarySkill:      "Go, Python",
					SecondarySkill:    "Docker, Kubernetes",
				}, nil).Once()
			},
			isErrorExpected: false,
			wantResponse: specs.IntranetEmployeeResponse{
				EmployeeID:        "EMP123",
				Email:             "emp@example.com",
				Name:              "John Doe",
				MobileNumber:      "1234567890",
				Gender:            "Male",
				YearsOfExperience: 5.5,
				Designation:       "Software Engineer",
				JoshJoiningDate:   "2020-01-01",
				LinkedinURL:       "https://linkedin.com/in/johndoe",
				GithubURL:         "https://github.com/johndoe",
				PrimarySkills:     []string{"Go", "Python"},
				SecondarySkills:   []string{"Docker", "Kubernetes"},
			},
			expectedError: nil,
		},
		{
			name:       "Fail_profile_already_exists",
			employeeID: "EMP123",
			setup: func(clientMock *intranetMocks.IntranetClient, repoMock *repoMocks.ProfileStorer) {
				repoMock.On("BeginTransaction", mock.Anything).Return(nil, nil)
				repoMock.On("GetProfileIDByEmployeeID", mock.Anything, "EMP123", mock.Anything).Return(1, nil).Once()
				repoMock.On("GetProfile", mock.Anything, 1, mock.Anything).Return(specs.ResponseProfile{Name: "Test User"}, nil).Once()
				repoMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			isErrorExpected: true,
			wantResponse:    specs.IntranetEmployeeResponse{},
			expectedError:   errs.ProfileExistsError{Name: "Test User"},
		},
		{
			name:       "Fail_intranet_client_returns_error",
			employeeID: "EMP123",
			setup: func(clientMock *intranetMocks.IntranetClient, repoMock *repoMocks.ProfileStorer) {
				repoMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				repoMock.On("GetProfileIDByEmployeeID", mock.Anything, "EMP123", mock.Anything).Return(0, errs.ErrNoRecordFound).Once()
				repoMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				clientMock.On("GetEmployeeByID", mock.Anything, "EMP123").Return(nil, errors.New("upstream error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    specs.IntranetEmployeeResponse{},
			expectedError:   errors.New("upstream error"),
		},
		{
			name:       "Fail_intranet_employee_not_found",
			employeeID: "EMP404",
			setup: func(clientMock *intranetMocks.IntranetClient, repoMock *repoMocks.ProfileStorer) {
				repoMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				repoMock.On("GetProfileIDByEmployeeID", mock.Anything, "EMP404", mock.Anything).Return(0, errs.ErrNoRecordFound).Once()
				repoMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				clientMock.On("GetEmployeeByID", mock.Anything, "EMP404").Return(nil, errs.ErrNoRecordFound).Once()
			},
			isErrorExpected: true,
			wantResponse:    specs.IntranetEmployeeResponse{},
			expectedError:   errs.ErrNoRecordFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset mocks for each test
			mockIntranetClient.ExpectedCalls = nil
			mockProfileRepo.ExpectedCalls = nil

			test.setup(mockIntranetClient, mockProfileRepo)
			gotResp, err := profileService.GetIntranetEmployee(context.Background(), test.employeeID)

			assert.Equal(t, test.wantResponse, gotResp)

			if test.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockIntranetClient.AssertExpectations(t)
			mockProfileRepo.AssertExpectations(t)
		})
	}
}
