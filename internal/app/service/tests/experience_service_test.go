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
