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
				Achievements: []dto.Achievement{
					{
						Name:        "Star Performer",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateAchievement returns an error",
			input: dto.CreateAchievementRequest{
				Achievements: []dto.Achievement{
					{
						Name:        "Achievement B",
						Description: "Description B",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of missing achievement name",
			input: dto.CreateAchievementRequest{
				Achievements: []dto.Achievement{
					{
						Name:        "",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo")).Return(errors.New("Missing achievement name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of empty payload",
			input: dto.CreateAchievementRequest{
				Achievements: []dto.Achievement{},
			},
			setup:           func(achievementMock *mocks.AchievementStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockAchievementRepo)

			_, err := achService.CreateAchievement(context.TODO(), test.input, "1")

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestUpdateAchievement(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	var repodeps = service.RepoDeps{
		AchievementDeps: mockAchievementRepo,
	}
	achService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       string
		achievementID   string
		input           dto.UpdateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success for updating achievement details",
			profileID:     "1",
			achievementID: "1",
			input: dto.UpdateAchievementRequest{
				Achievement: dto.Achievement{
					Name:        "Updated Star Performer",
					Description: "Updated Description",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("UpdateAchievement", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed because UpdateAchievement returns an error",
			profileID:     "100000000000000000",
			achievementID: "1",
			input: dto.UpdateAchievementRequest{
				Achievement: dto.Achievement{
					Name:        "Achievement B",
					Description: "Description B",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("UpdateAchievement", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateAchievementRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed because of missing achievement name",
			profileID:     "1",
			achievementID: "1",
			input: dto.UpdateAchievementRequest{
				Achievement: dto.Achievement{
					Name:        "",
					Description: "Description",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer) {
				achievementMock.On("UpdateAchievement", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo")).Return(0, errors.New("Missing achievement name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed because of invalid profileID or achievementID",
			profileID:     "invalid",
			achievementID: "1",
			input: dto.UpdateAchievementRequest{
				Achievement: dto.Achievement{
					Name:        "Valid Name",
					Description: "Valid Description",
				},
			},
			setup:           func(achievementMock *mocks.AchievementStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockAchievementRepo)

			_, err := achService.UpdateAchievement(context.TODO(), test.profileID, test.achievementID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
