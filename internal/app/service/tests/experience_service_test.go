package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockResponseExperience = []dto.ExperienceResponse{
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
		input           dto.CreateExperienceRequest
		setup           func(experienceMock *mocks.ExperienceStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for experience details",
			input: dto.CreateExperienceRequest{
				ProfileID: 1,
				Experiences: []dto.Experience{
					{
						Designation: "Software Engineer",
						CompanyName: "Josh Software Pvt.Ltd.",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateExperience returns an error",
			input: dto.CreateExperienceRequest{
				ProfileID: 10000000000,
				Experiences: []dto.Experience{
					{
						Designation: "Software Engineer",
						CompanyName: "Tech Corp",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of missing designation",
			input: dto.CreateExperienceRequest{
				ProfileID: 1,
				Experiences: []dto.Experience{
					{
						Designation: "",
						CompanyName: "Tech Corp",
						FromDate:    "2023-01-01",
						ToDate:      "2024-01-01",
					},
				},
			},
			setup: func(experienceMock *mocks.ExperienceStorer) {
				experienceMock.On("CreateExperience", mock.Anything, mock.AnythingOfType("[]repository.ExperienceDao")).Return(errors.New("Missing designation")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of empty payload",
			input: dto.CreateExperienceRequest{
				ProfileID:   1,
				Experiences: []dto.Experience{},
			},
			setup:           func(experienceMock *mocks.ExperienceStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockExperienceRepo)

			// Test the service
			_, err := experienceService.CreateExperience(context.TODO(), test.input)

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
		profileID       string
		setup           func(expMock *mocks.ExperienceStorer)
		isErrorExpected bool
		wantResponse    []dto.ExperienceResponse
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
				expMock.On("GetExperiences", mock.Anything, mock.Anything).Return([]dto.ExperienceResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []dto.ExperienceResponse{},
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
