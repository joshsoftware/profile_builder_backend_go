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

var mockListEduFilter = specs.ListEducationsFilter{
	EduationsIDs: []int{},
	Names:        []string{},
}

func TestCreateEducation(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	var repodeps = service.RepoDeps{
		EducationDeps: mockEducationRepo,
	}
	eduService := service.NewServices(repodeps)

	tests := []struct {
		name              string
		input             specs.CreateEducationRequest
		setup             func(educationMock *mocks.EducationStorer)
		isErrorExpected   bool
		expectedProfileID int
	}{
		{
			name: "Success_for_education_details",
			input: specs.CreateEducationRequest{
				Educations: []specs.Education{
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
			name: "Failed_because_of_error",
			input: specs.CreateEducationRequest{
				Educations: []specs.Education{
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
			name: "Failed_because_of_empty_payload",
			input: specs.CreateEducationRequest{
				Educations: []specs.Education{},
			},
			setup:             func(educationMock *mocks.EducationStorer) {},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo)

			profileID, err := eduService.CreateEducation(context.Background(), test.input, 1, 1)
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

func TestListEducations(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	var repodeps = service.RepoDeps{
		EducationDeps: mockEducationRepo,
	}
	educationService := service.NewServices(repodeps)

	mockProfileID := 123
	mockResponseEducation := []specs.EducationResponse{
		{
			ProfileID:        123,
			Degree:           "B.Sc. Computer Science",
			UniversityName:   "University of Example",
			Place:            "Example City",
			PercentageOrCgpa: "3.8",
			PassingYear:      "2020",
		},
	}

	tests := []struct {
		name            string
		profileID       int
		setup           func(eduMock *mocks.EducationStorer)
		isErrorExpected bool
		wantResponse    []specs.EducationResponse
	}{
		{
			name:      "Success_get_education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer) {
				// Mock successful retrieval
				eduMock.On("GetEducation", mock.Anything, mock.Anything).Return(mockResponseEducation, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseEducation,
		},
		{
			name:      "Fail_get_education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer) {
				// Mock retrieval failure
				eduMock.On("GetEducation", mock.Anything, mock.Anything).Return([]specs.EducationResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.EducationResponse{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockEducationRepo)

			// Call the method being tested
			gotResp, err := educationService.ListEducations(context.Background(), test.profileID, mockListEduFilter)

			// Assertions
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
		})
	}
}

func TestUpdateEducation(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	var repodeps = service.RepoDeps{
		EducationDeps: mockEducationRepo,
	}
	eduService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		educationID     int
		userID          int
		input           specs.UpdateEducationRequest
		setup           func(educationMock *mocks.EducationStorer)
		isErrorExpected bool
	}{
		{
			name:        "Success_for_updating_education_details",
			profileID:   1,
			educationID: 1,
			userID:      1,
			input: specs.UpdateEducationRequest{
				Education: specs.Education{
					Degree:           "Updated Degree",
					UniversityName:   "Updated University",
					Place:            "Updated Place",
					PercentageOrCgpa: "4.0",
					PassingYear:      "2023",
				},
			},
			setup: func(educationMock *mocks.EducationStorer) {
				educationMock.On("UpdateEducation", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateEducationRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:        "Failed_because_updateeducation_returns_an_error",
			profileID:   100000,
			educationID: 1,
			userID:      1,
			input: specs.UpdateEducationRequest{
				Education: specs.Education{
					Degree:           "Degree B",
					UniversityName:   "University B",
					Place:            "Place B",
					PercentageOrCgpa: "3.5",
					PassingYear:      "2022",
				},
			},
			setup: func(educationMock *mocks.EducationStorer) {
				educationMock.On("UpdateEducation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateEducationRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:        "Failed_because_of_missing_education_degree",
			profileID:   1,
			educationID: 1,
			userID:      1,
			input: specs.UpdateEducationRequest{
				Education: specs.Education{
					Degree:           "",
					UniversityName:   "University",
					Place:            "Place",
					PercentageOrCgpa: "3.8",
					PassingYear:      "2023",
				},
			},
			setup: func(educationMock *mocks.EducationStorer) {
				educationMock.On("UpdateEducation", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateEducationRepo")).Return(0, errors.New("Missing education degree")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:        "Failed_because_of_invalid_profileid_or_educationid",
			profileID:   -1,
			educationID: 1,
			userID:      1,
			input: specs.UpdateEducationRequest{
				Education: specs.Education{
					Degree:           "Valid Degree",
					UniversityName:   "Valid University",
					Place:            "Valid Place",
					PercentageOrCgpa: "3.9",
					PassingYear:      "2023",
				},
			},
			setup:           func(educationMock *mocks.EducationStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo)

			_, err := eduService.UpdateEducation(context.TODO(), test.profileID, test.educationID, test.userID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
