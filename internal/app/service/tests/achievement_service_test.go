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

func TestCreateAchievement(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	var repodeps = service.RepoDeps{
		AchievementDeps: mockAchievementRepo,
	}
	achService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           dto.CreateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for achievement details",
			input: dto.CreateAchievementRequest{
				ProfileID: 1,
				Achievements: []dto.Achievement{
					{
						Name:        "Star Performer",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateAchievement returns an error",
			input: dto.CreateAchievementRequest{
				ProfileID: 100000000000000000,
				Achievements: []dto.Achievement{
					{
						Name:        "Achievement B",
						Description: "Description B",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of missing achievement name",
			input: dto.CreateAchievementRequest{
				ProfileID: 1,
				Achievements: []dto.Achievement{
					{
						Name:        "",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementDao")).Return(errors.New("Missing achievement name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of empty payload",
			input: dto.CreateAchievementRequest{
				ProfileID:    1,
				Achievements: []dto.Achievement{},
			},
			setup:           func(achievementMock *mocks.AchievementStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockAchievementRepo)

			_, err := achService.CreateAchievement(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
