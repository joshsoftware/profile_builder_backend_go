package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
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
		input           specs.CreateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_for_achievement_details",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{
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
			name: "Failed_because_createachievement_returns_an_error",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{
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
			name: "Failed_because_of_missing_achievement_name",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{
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
			name: "Failed_because_of_empty_payload",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{},
			},
			setup:           func(achievementMock *mocks.AchievementStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockAchievementRepo)

			_, err := achService.CreateAchievement(context.TODO(), test.input, 1)

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
		input           specs.UpdateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_updating_achievement_details",
			profileID:     "1",
			achievementID: "1",
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
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
			name:          "Failed_because_updateachievement_returns_an_error",
			profileID:     "100000000000000000",
			achievementID: "1",
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
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
			name:          "Failed_because_of_missing_achievement_name",
			profileID:     "1",
			achievementID: "1",
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
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
			name:          "Failed_because_of_invalid_profileid_or_achievementID",
			profileID:     "invalid",
			achievementID: "1",
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
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

var (
	profileID = 123
)

func TestListAchievements(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	var repodeps = service.RepoDeps{
		AchievementDeps: mockAchievementRepo,
	}
	achService := service.NewServices(repodeps)

	// mock data
	mockProfileId := profileID
	mockResponseAchievement := []specs.AchievementResponse{
		{
			ProfileID:   123,
			Name:        "Client Appreciation",
			Description: "Appreciated by client for the work done",
		},
	}

	tests := []struct {
		Name            string
		ProfileID       string
		MockSetup       func(*mocks.AchievementStorer, int)
		isErrorExpected bool
		wantResponse    []specs.AchievementResponse
	}{
		{
			Name:      "success_get_achievement",
			ProfileID: strconv.Itoa(mockProfileId),
			MockSetup: func(mockAchievementStorer *mocks.AchievementStorer, profileID int) {
				mockAchievementStorer.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return(mockResponseAchievement, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseAchievement,
		},
		{
			Name:      "fail_get_achievement",
			ProfileID: "123",
			MockSetup: func(achMock *mocks.AchievementStorer, profileID int) {
				achMock.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.AchievementResponse{},
		},
		{
			Name:      "sucess_with_empty_resultset",
			ProfileID: "123",
			MockSetup: func(achMock *mocks.AchievementStorer, profileID int) {
				achMock.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{}, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    []specs.AchievementResponse{},
		},
		{
			Name:      "invalid_profile_id",
			ProfileID: "invalid",
			MockSetup: func(achMock *mocks.AchievementStorer, profileID int) {
				achMock.On("ListAchievements", mock.Anything, mock.Anything, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("invalid profile ID")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.AchievementResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			profileIDInt, _ := strconv.Atoi(tt.ProfileID)
			tt.MockSetup(mockAchievementRepo, profileIDInt)
			gotResponse, err := achService.ListAchievements(context.Background(), profileIDInt, specs.ListAchievementFilter{})

			assert.Equal(t, tt.wantResponse, gotResponse)
			if (err != nil) != tt.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", tt.Name, tt.isErrorExpected, err != nil)
			}
		})
	}
}
