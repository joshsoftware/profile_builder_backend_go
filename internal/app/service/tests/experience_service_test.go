package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockListExpFilter = specs.ListExperiencesFilter{
	ExperiencesIDs: []int{},
	Names:          []string{},
}

var mockResponseExperience = []specs.ExperienceResponse{
	{
		ProfileID:   123,
		Designation: "Software Engineer",
		CompanyName: "Tech Corp",
		FromDate:    "2018-01-01",
		ToDate:      "2020-12-31",
	},
}

func TestCreateExperience(t *testing.T) {
	mockExperienceRepo := new(mocks.ExperienceStorer)
	var repodeps = service.RepoDeps{
		ExperienceDeps: mockExperienceRepo,
	}
	experienceService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           specs.CreateExperienceRequest
		setup           func(experienceMock *mocks.ExperienceStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_for_experience_details",
			input: specs.CreateExperienceRequest{
				Experiences: []specs.Experience{
					{
						Designation: "Software Engineer",
						CompanyName: "Josh Software Pvt.Ltd.",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceRepo")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed_because_createexperience_returns_an_error",
			input: specs.CreateExperienceRequest{
				Experiences: []specs.Experience{
					{
						Designation: "Software Engineer",
						CompanyName: "Tech Corp",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceRepo")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_missing_designation",
			input: specs.CreateExperienceRequest{
				Experiences: []specs.Experience{
					{
						Designation: "",
						CompanyName: "Tech Corp",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceRepo")).Return(errors.New("Missing designation")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_empty_payload",
			input: specs.CreateExperienceRequest{
				Experiences: []specs.Experience{},
			},
			setup:           func(experienceMock *mocks.ExperienceStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockExperienceRepo)

			// Test the service
			_, err := experienceService.CreateExperience(context.TODO(), test.input, 1, 1)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestGetExperience(t *testing.T) {
	// Initialize mock dependencies
	mockExperienceRepo := new(mocks.ExperienceStorer)
	var repodeps = service.RepoDeps{
		ExperienceDeps: mockExperienceRepo,
	}
	experienceService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		setup           func(expMock *mocks.ExperienceStorer)
		isErrorExpected bool
		wantResponse    []specs.ExperienceResponse
	}{
		{
			name:      "Success_get_experience",
			profileID: mockProfileID,
			setup: func(expMock *mocks.ExperienceStorer) {
				expMock.On("ListExperiences", mock.Anything, mock.Anything).Return(mockResponseExperience, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseExperience,
		},
		{
			name:      "Fail_get_experience",
			profileID: mockProfileID,
			setup: func(expMock *mocks.ExperienceStorer) {
				// Mock retrieval failure
				expMock.On("ListExperiences", mock.Anything, mock.Anything).Return([]specs.ExperienceResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.ExperienceResponse{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockExperienceRepo)

			// Call the method being tested
			gotResp, err := experienceService.ListExperiences(context.Background(), test.profileID, mockListExpFilter)

			// Assertions
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
		})
	}
}

func TestUpdateExperience(t *testing.T) {
	mockExperienceRepo := new(mocks.ExperienceStorer)
	var repodeps = service.RepoDeps{
		ExperienceDeps: mockExperienceRepo,
	}
	expService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		experienceID    int
		userID          int
		input           specs.UpdateExperienceRequest
		setup           func(experienceMock *mocks.ExperienceStorer)
		isErrorExpected bool
	}{
		{
			name:         "Success_for_updating_experience_details",
			profileID:    1,
			experienceID: 1,
			userID:       1,
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "Updated Designation",
					CompanyName: "Updated Company",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:         "Failed_because_updateexperience_returns_an_error",
			profileID:    100000,
			experienceID: 1,
			userID:       1,
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "Designation B",
					CompanyName: "Company B",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:         "Failed_because_of_missing_experience_designation",
			profileID:    1,
			experienceID: 1,
			userID:       1,
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "",
					CompanyName: "Company",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(0, errors.New("Missing experience designation")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:         "Failed_because_of_invalid_profileid_or_experienceid",
			profileID:    -1,
			experienceID: 1,
			userID:       1,
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "Valid Designation",
					CompanyName: "Valid Company",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup:           func(experienceMock *mocks.ExperienceStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockExperienceRepo)

			_, err := expService.UpdateExperience(context.TODO(), test.profileID, test.experienceID, test.userID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestDeleteExperienceService(t *testing.T) {
	mockExperienceSvc := new(mocks.ExperienceStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		ExperienceDeps: mockExperienceSvc,
		ProfileDeps:    mockProfileRepo,
	}
	experienceSvc := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		experienceID    int
		profileID       int
		setup           func(experienceMock *mocks.ExperienceStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:         "Success_for_delete_experience",
			experienceID: 1,
			profileID:    1,
			setup: func(experienceMock *mocks.ExperienceStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				experienceMock.On("DeleteExperience", mock.Anything, 1, 1, nil).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:         "Failed_because_delete_experience_returns_an_error",
			experienceID: 2,
			profileID:    1,
			setup: func(experienceMock *mocks.ExperienceStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				experienceMock.On("DeleteExperience", mock.Anything, 1, 2, nil).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:         "Failed_because_DeleteExperience_returns_an_error",
			experienceID: 3,
			profileID:    1,
			setup: func(experienceMock *mocks.ExperienceStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				experienceMock.On("DeleteExperience", mock.Anything, 1, 3, nil).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, nil, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		// {
		// 	name: "Failed_because_BeginTransaction_returns_an_error",
		// 		achievementID: 4,
		// 		profileID:     1,
		// 	setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
		// 		profileMock.On("BeginTransaction", mock.Anything).Return(nil, errors.New("error")).Once()
		// 	},
		// 	isErrorExpected: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockExperienceSvc, mockProfileRepo)
			err := experienceSvc.DeleteExperience(context.Background(), test.profileID, test.experienceID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
