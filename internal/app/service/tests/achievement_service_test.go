package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getMock(t *testing.T) *mocks.AchievementStorer {
	mock := &mocks.AchievementStorer{}
	mock.Test(t)
	return mock
}

func getProfileMock(t *testing.T) *mocks.ProfileStorer {
	mock := &mocks.ProfileStorer{}
	mock.Test(t)
	return mock
}

func TestCreateAchievement(t *testing.T) {
	mockAchievementRepo := getMock(t)
	mockProfileRepo := getProfileMock(t)

	repoDeps := service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
		AchievementDeps: mockAchievementRepo,
	}
	achService := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		input           specs.CreateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer)
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
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, nil).Once()
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo"), mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
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
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, nil).Once()
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo"), mock.Anything).Return(errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
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
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, nil).Once()
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo"), mock.Anything).Return(errors.New("missing achievement name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("transaction handling failed")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_empty_payload",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, nil).Once()
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo"), mock.Anything).Return(errors.New("empty payload")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("transaction handling failed")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_counting_records_fails",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{
					{
						Name:        "Star Performer",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, errors.New("Error counting records")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_transaction_handling_failed",
			input: specs.CreateAchievementRequest{
				Achievements: []specs.Achievement{
					{
						Name:        "Star Performer",
						Description: "Description",
					},
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, constants.Achievements, mock.Anything).Return(0, nil).Once()
				achievementMock.On("CreateAchievement", mock.Anything, mock.AnythingOfType("[]repository.AchievementRepo"), mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("transaction handling failed")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup(mockAchievementRepo, mockProfileRepo)
			}
			_, err := achService.CreateAchievement(context.TODO(), test.input, 1, 1)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
			mockProfileRepo.AssertExpectations(t)
			mockAchievementRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateAchievement(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	mockProfileRepo := new(mocks.ProfileStorer)

	repodeps := service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
		AchievementDeps: mockAchievementRepo,
	}
	achService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		achievementID   int
		userID          int
		input           specs.UpdateAchievementRequest
		setup           func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_updating_achievement_details",
			profileID:     1,
			achievementID: 1,
			userID:        1,
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
					Name:        "Updated Star Performer",
					Description: "Updated Description",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("UpdateAchievement", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed_because_updateachievement_returns_an_error",
			profileID:     1,
			achievementID: 1,
			userID:        1,
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
					Name:        "Achievement B",
					Description: "Description B",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("UpdateAchievement", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error in handle transaction")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_missing_achievement_name",
			profileID:     1,
			achievementID: 1,
			userID:        1,
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
					Name:        "",
					Description: "Description",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("UpdateAchievement", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo"), mock.Anything).Return(0, errors.New("Missing achievement name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("transaction handling failed")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_invalid_profileid_or_achievementID",
			profileID:     -1,
			achievementID: 1,
			userID:        1,
			input: specs.UpdateAchievementRequest{
				Achievement: specs.Achievement{
					Name:        "Valid Name",
					Description: "Valid Description",
				},
			},
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("UpdateAchievement", mock.Anything, -1, 1, mock.AnythingOfType("repository.UpdateAchievementRepo"), mock.Anything).Return(0, errors.New("Invalid profileID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("transaction handling failed")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup(mockAchievementRepo, mockProfileRepo)
			}
			_, err := achService.UpdateAchievement(context.TODO(), test.profileID, test.achievementID, test.userID, test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
			mockProfileRepo.AssertExpectations(t)
			mockAchievementRepo.AssertExpectations(t)
		})
	}
}

var (
	profileID = 123
)

func TestListAchievements(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
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
		MockSetup       func(*mocks.AchievementStorer, *mocks.ProfileStorer, int)
		isErrorExpected bool
		wantResponse    []specs.AchievementResponse
	}{
		{
			Name:      "success_get_achievement",
			ProfileID: strconv.Itoa(mockProfileId),
			MockSetup: func(mockAchievementStorer *mocks.AchievementStorer, mockProfileStorer *mocks.ProfileStorer, profileID int) {
				mockProfileStorer.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockAchievementStorer.On("ListAchievements", mock.Anything, profileID, mock.Anything, mock.Anything).Return(mockResponseAchievement, nil).Once()
				mockProfileStorer.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseAchievement,
		},
		{
			Name:      "fail_get_achievement",
			ProfileID: "123",
			MockSetup: func(achMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achMock.On("ListAchievements", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.AchievementResponse{},
		},
		{
			Name:      "sucess_with_empty_resultset",
			ProfileID: "123",
			MockSetup: func(achMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achMock.On("ListAchievements", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.AchievementResponse{}, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    []specs.AchievementResponse{},
		},
		{
			Name:      "invalid_profile_id",
			ProfileID: "invalid",
			MockSetup: func(achMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achMock.On("ListAchievements", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("invalid profile ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.AchievementResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			profileIDInt, _ := strconv.Atoi(tt.ProfileID)
			tt.MockSetup(mockAchievementRepo, mockProfileRepo, profileIDInt)
			gotResponse, err := achService.ListAchievements(context.Background(), profileIDInt, specs.ListAchievementFilter{})

			assert.Equal(t, tt.wantResponse, gotResponse)
			if (err != nil) != tt.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", tt.Name, tt.isErrorExpected, err != nil)
			}
		})
	}
}

func TestDeleteAchievementService(t *testing.T) {
	mockAchievementRepo := new(mocks.AchievementStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		AchievementDeps: mockAchievementRepo,
		ProfileDeps:     mockProfileRepo,
	}
	achService := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		achievementID   int
		profileID       int
		setup           func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_delete_achievement",
			achievementID: 1,
			profileID:     1,
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("DeleteAchievement", mock.Anything, 1, 1, nil).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed_because_delete_achievement_returns_an_error",
			achievementID: 2,
			profileID:     1,
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("DeleteAchievement", mock.Anything, 1, 2, nil).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_DeleteAchievement_returns_an_error",
			achievementID: 3,
			profileID:     1,
			setup: func(achievementMock *mocks.AchievementStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				achievementMock.On("DeleteAchievement", mock.Anything, 1, 3, nil).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, nil, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockAchievementRepo, mockProfileRepo)
			err := achService.DeleteAchievement(context.Background(), test.profileID, test.achievementID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}

			mockAchievementRepo.AssertExpectations(t)
			mockProfileRepo.AssertExpectations(t)
		})
	}
}
