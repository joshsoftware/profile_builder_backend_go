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

var mockListEduFilter = specs.ListEducationsFilter{
	EduationsIDs: []int{},
	Names:        []string{},
}

func TestCreateEducation(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:   mockProfileRepo,
		EducationDeps: mockEducationRepo,
	}
	eduService := service.NewServices(repodeps)

	tests := []struct {
		name              string
		input             specs.CreateEducationRequest
		setup             func(*mocks.EducationStorer, *mocks.ProfileStorer)
		profileID         int
		userID            int
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
			profileID: 1,
			userID:    1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(2, nil).Once()
				educationMock.On("CreateEducation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected:   false,
			expectedProfileID: 1,
		},
		{
			name: "Success_for_multiple_education_entries",
			input: specs.CreateEducationRequest{
				Educations: []specs.Education{
					{
						Degree:           "B.Sc in Mathematics",
						UniversityName:   "XYZ University",
						Place:            "City A",
						PercentageOrCgpa: "3.7",
						PassingYear:      "2023",
					},
					{
						Degree:           "M.Sc in Mathematics",
						UniversityName:   "XYZ University",
						Place:            "City A",
						PercentageOrCgpa: "3.8",
						PassingYear:      "2025",
					},
				},
			},
			profileID: 1,
			userID:    1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(2, nil).Once()
				educationMock.On("CreateEducation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
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
			profileID: 1,
			userID:    1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				educationMock.On("CreateEducation", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
		{
			name: "Failed_because_of_empty_payload",
			input: specs.CreateEducationRequest{
				Educations: []specs.Education{},
			},
			profileID: 1,
			userID:    1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				educationMock.On("CreateEducation", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handler transaction error")).Once()
			},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
		{
			name: "Failed_due_to_record_count_error",
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
			profileID: 1,
			userID:    1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(0, errors.New("Count records error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handler transaction error")).Once()
			},
			isErrorExpected:   true,
			expectedProfileID: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo, mockProfileRepo)

			profileID, err := eduService.CreateEducation(context.Background(), test.input, 1, 1)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			if profileID != test.expectedProfileID {
				t.Errorf("Test %s failed, expected profileID to be %v, but got %v", test.name, test.expectedProfileID, profileID)
			}

			mockProfileRepo.AssertExpectations(t)
			mockEducationRepo.AssertExpectations(t)
		})
	}
}

func TestListEducations(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:   mockProfileRepo,
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
		setup           func(*mocks.EducationStorer, *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    []specs.EducationResponse
	}{
		{
			name:      "Success_get_education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				eduMock.On("ListEducations", mock.Anything, mockProfileID, mock.Anything, mock.Anything).Return(mockResponseEducation, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseEducation,
		},
		{
			name:      "Fail_get_education",
			profileID: mockProfileID,
			setup: func(eduMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				eduMock.On("ListEducations", mock.Anything, mockProfileID, mock.Anything, mock.Anything).Return([]specs.EducationResponse{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.EducationResponse{},
		},
		{
			name:      "Fail_get_education_due_to_invalid_profile_id",
			profileID: -1,
			setup: func(eduMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				eduMock.On("ListEducations", mock.Anything, -1, mock.Anything, mock.Anything).Return([]specs.EducationResponse{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.EducationResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo, mockProfileRepo)
			gotResp, err := educationService.ListEducations(context.Background(), test.profileID, mockListEduFilter)
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}

			mockProfileRepo.AssertExpectations(t)
			mockEducationRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateEducation(t *testing.T) {
	mockEducationRepo := new(mocks.EducationStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:   mockProfileRepo,
		EducationDeps: mockEducationRepo,
	}
	eduService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		educationID     int
		userID          int
		input           specs.UpdateEducationRequest
		setup           func(*mocks.EducationStorer, *mocks.ProfileStorer)
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
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("UpdateEducation", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateEducationRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:        "Failed_because_updateeducation_returns_an_error",
			profileID:   1,
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
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("UpdateEducation", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateEducationRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
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
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("UpdateEducation", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateEducationRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
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
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("UpdateEducation", mock.Anything, -1, 1, mock.AnythingOfType("repository.UpdateEducationRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationRepo, mockProfileRepo)
			_, err := eduService.UpdateEducation(context.TODO(), test.profileID, test.educationID, test.userID, test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
			mockEducationRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteEducationService(t *testing.T) {
	mockEducationSvc := new(mocks.EducationStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		EducationDeps: mockEducationSvc,
		ProfileDeps:   mockProfileRepo,
	}
	educationSvc := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		educationID     int
		profileID       int
		setup           func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:        "Success_for_delete_education",
			educationID: 1,
			profileID:   1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("DeleteEducation", mock.Anything, 1, 1, nil).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:        "Failed_because_delete_education_returns_an_error",
			educationID: 2,
			profileID:   1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("DeleteEducation", mock.Anything, 1, 2, nil).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:        "Failed_because_DeleteEducation_returns_an_error",
			educationID: 3,
			profileID:   1,
			setup: func(educationMock *mocks.EducationStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				educationMock.On("DeleteEducation", mock.Anything, 1, 3, nil).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, nil, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockEducationSvc, mockProfileRepo)
			err := educationSvc.DeleteEducation(context.Background(), test.profileID, test.educationID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}

			mockProfileRepo.AssertExpectations(t)
			mockEducationSvc.AssertExpectations(t)
		})
	}
}
