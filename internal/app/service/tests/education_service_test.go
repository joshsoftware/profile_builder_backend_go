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

func TestCreateEducation(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	var repodeps = service.RepoDeps{
		EducationDeps: mockEducationRepo,
	}
	eduService := service.NewServices(repodeps)

	tests := []struct {
		name              string
		input             dto.CreateEducationRequest
		setup             func(educationMock *mocks.EducationStorer)
		isErrorExpected   bool
		expectedProfileID int
	}{
		{
			name: "Success for education details",
			input: dto.CreateEducationRequest{
				ProfileID: 1,
				Educations: []dto.Education{
					{
						Degree:           "B.Tech in Computer Science",
						UniversityName:   "SPPU",
						Place:            "Pune",
						PercentageOrCgpa: "3.5",
						PassingYear:      "2022",
					},
				},
			},
			setup: func(educationMock *mocks.EducationStorer) {
				educationMock.On("CreateEducation", mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected:   false,
			expectedProfileID: 1,
		},
		{
			name: "Failed because of error",
			input: dto.CreateEducationRequest{
				ProfileID: 456,
				Educations: []dto.Education{
					{
						Degree:           "",
						UniversityName:   "Example University",
						Place:            "Example City",
						PercentageOrCgpa: "3.5",
						PassingYear:      "2022",
					},
				},
			},
			setup: func(educationMock *mocks.EducationStorer) {
				educationMock.On("CreateEducation", mock.Anything, mock.Anything).Return(errors.New("Error")).Once()
			},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
		{
			name: "Failed because of empty payload",
			input: dto.CreateEducationRequest{
				ProfileID:  456,
				Educations: []dto.Education{},
			},
			setup:             func(educationMock *mocks.EducationStorer) {},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo)

			profileID, err := eduService.CreateEducation(context.Background(), test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			if profileID != test.expectedProfileID {
				t.Errorf("Test %s failed, expected profileID to be %v, but got %v", test.name, test.expectedProfileID, profileID)
			}

			mockEducationRepo.AssertExpectations(t)
		})
	}
}

func TestGetEducation(t *testing.T) {
	// Initialize mock dependencies
	mockEducationRepo := new(mocks.EducationStorer)
	var repodeps = service.RepoDeps{
		EducationDeps: mockEducationRepo,
	}
	educationService := service.NewServices(repodeps)

	// Define mock data
	mockProfileID := "123"
	mockResponseEducation := []dto.EducationResponse{
		{
			ProfileID:        123,
			Degree:           "B.Sc. Computer Science",
			UniversityName:   "University of Example",
			Place:            "Example City",
			PercentageOrCgpa: "3.8",
			PassingYear:      "2020",
		},
	}

	// Define test cases
	tests := []struct {
		name            string
		profileID       string
		setup           func(eduMock *mocks.EducationStorer)
		isErrorExpected bool
		wantResponse    []dto.EducationResponse
	}{
		{
			name:      "Success get education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer) {
				// Mock successful retrieval
				eduMock.On("GetEducation", mock.Anything, mock.Anything).Return(mockResponseEducation, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseEducation,
		},
		{
			name:      "Fail get education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer) {
				// Mock retrieval failure
				eduMock.On("GetEducation", mock.Anything, mock.Anything).Return([]dto.EducationResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []dto.EducationResponse{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockEducationRepo)

			// Call the method being tested
			gotResp, err := educationService.GetEducation(context.Background(), test.profileID)

			// Assertions
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
		})
	}
}
