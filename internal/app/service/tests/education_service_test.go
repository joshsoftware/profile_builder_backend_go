package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
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
