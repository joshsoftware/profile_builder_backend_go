package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
			name: "Success for experience details",
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
			name: "Failed because CreateExperience returns an error",
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
			name: "Failed because of missing designation",
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
			name: "Failed because of empty payload",
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
			_, err := experienceService.CreateExperience(context.TODO(), test.input, 1)

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
			name:      "Success get experience",
			profileID: mockProfileID,
			setup: func(expMock *mocks.ExperienceStorer) {
				// Mock successful retrieval
				expMock.On("GetExperiences", mock.Anything, mock.Anything).Return(mockResponseExperience, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseExperience,
		},
		{
			name:      "Fail get experience",
			profileID: mockProfileID,
			setup: func(expMock *mocks.ExperienceStorer) {
				// Mock retrieval failure
				expMock.On("GetExperiences", mock.Anything, mock.Anything).Return([]specs.ExperienceResponse{}, errors.New("error")).Once()
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
			gotResp, err := experienceService.GetExperience(context.Background(), test.profileID)

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
		profileID       string
		experienceID    string
		input           specs.UpdateExperienceRequest
		setup           func(experienceMock *mocks.ExperienceStorer)
		isErrorExpected bool
	}{
		{
			name:         "Success for updating experience details",
			profileID:    "1",
			experienceID: "1",
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "Updated Designation",
					CompanyName: "Updated Company",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:         "Failed because UpdateExperience returns an error",
			profileID:    "100000000000000000",
			experienceID: "1",
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "Designation B",
					CompanyName: "Company B",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:         "Failed because of missing experience designation",
			profileID:    "1",
			experienceID: "1",
			input: specs.UpdateExperienceRequest{
				Experience: specs.Experience{
					Designation: "",
					CompanyName: "Company",
					FromDate:    "2022-01-01",
					ToDate:      "2023-01-01",
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("UpdateExperience", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateExperienceRepo")).Return(0, errors.New("Missing experience designation")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:         "Failed because of invalid profileID or experienceID",
			profileID:    "invalid",
			experienceID: "1",
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

			_, err := expService.UpdateExperience(context.TODO(), test.profileID, test.experienceID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
