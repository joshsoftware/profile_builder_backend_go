package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockProfile = specs.Profile{
	Name:              "Example User",
	Email:             "example.user@gmail.com",
	Gender:            "Male",
	Mobile:            "9955662233",
	Designation:       "Software Engineer",
	Description:       "Experienced software engineer",
	Title:             "Golang Developer",
	YearsOfExperience: 5,
	PrimarySkills:     []string{"Golang", "Python"},
	SecondarySkills:   []string{"JavaScript", "SQL"},
	GithubLink:        "https://github.com/abhishek",
	LinkedinLink:      "https://www.linkedin.com/in/abhsihek",
}

var mockProfileID = 123

var mockResponseProfile = specs.ResponseProfile{
	ProfileID:         123,
	Name:              "Example User",
	Email:             "example.user@gmail.com",
	Gender:            "Male",
	Mobile:            "9955662233",
	Designation:       "Software Engineer",
	Description:       "Experienced software engineer",
	Title:             "Golang Developer",
	YearsOfExperience: 5,
	PrimarySkills:     []string{"Golang", "Python"},
	SecondarySkills:   []string{"JavaScript", "SQL"},
	GithubLink:        "https://github.com/abhishek",
	LinkedinLink:      "https://www.linkedin.com/in/abhishek",
}

var mockListSkills = specs.ListSkills{Name: []string{"GO", "RUBY", "C", "C++", "JAVA", "PYTHON", "JAVASCRIPT"}}

func TestListProfile(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	mockListProfile := []specs.ListProfiles{
		{
			ID:                1,
			Name:              "Example profile",
			Email:             "example.profile@gmail.com",
			YearsOfExperience: 1.0,
			PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
			IsCurrentEmployee: 1,
			IsActive:          1,
			JoshJoiningDate:   sql.NullString{String: "2021-01-01", Valid: true},
			CreatedAt:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			IsProfileComplete: 1,
		},
	}

	mockResponseListProfile := []specs.ResponseListProfiles{
		{
			ID:                1,
			Name:              "Example profile",
			Email:             "example.profile@gmail.com",
			YearsOfExperience: 1.0,
			PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
			IsCurrentEmployee: "YES",
			IsActive:          "YES",
			JoshJoiningDate:   sql.NullString{String: "2021-01-01", Valid: true},
			CreatedAt:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			IsProfileComplete: "YES",
		},
	}

	tests := []struct {
		name            string
		setup           func(userMock *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    []specs.ResponseListProfiles
	}{
		{
			name: "Success_get_list_of_Profiles",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				userMock.On("ListProfiles", mock.Anything, mock.Anything).Return(mockListProfile, nil).Once()
				userMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseListProfile,
		},
		{
			name: "Fail_get_list_of_Profiles",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				userMock.On("ListProfiles", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
				userMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.ResponseListProfiles{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			gotResp, err := profileService.ListProfiles(context.Background())
			assert.Equal(t, test.wantResponse, gotResp)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestListSkills(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		setup           func(userMock *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    specs.ListSkills
	}{
		{
			name: "Success_get_list_of_Skills",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				userMock.On("ListSkills", mock.Anything, mock.Anything).Return(mockListSkills, nil).Once()
				userMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockListSkills,
		},
		{
			name: "Fail_get_list_of_Skills",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				userMock.On("ListSkills", mock.Anything, mock.Anything).Return(specs.ListSkills{}, errors.New("error")).Once()
				userMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    specs.ListSkills{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// Test service
			gotResp, err := profileService.ListSkills(context.Background())

			assert.Equal(t, test.wantResponse, gotResp)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateProfile(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           specs.CreateProfileRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_for_user_Detail",
			input: specs.CreateProfileRequest{
				Profile: mockProfile,
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("ProfileRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed_to_create_profile",
			input: specs.CreateProfileRequest{
				Profile: mockProfile,
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("ProfileRepo"), mock.Anything).Return(0, errors.New("profile creation failed")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			_, err := profileService.CreateProfile(context.TODO(), test.input, 1)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestGetProfile(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    specs.ResponseProfile
	}{
		{
			name:      "Success_get_profile",
			profileID: mockProfileID,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("GetProfile", mock.Anything, mock.Anything, mock.Anything).Return(mockResponseProfile, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseProfile,
		},
		{
			name:      "Fail_get_profile",
			profileID: mockProfileID,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("GetProfile", mock.Anything, mock.Anything, mock.Anything).Return(specs.ResponseProfile{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    specs.ResponseProfile{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			gotResp, err := profileService.GetProfile(context.Background(), test.profileID)
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		userID          int
		input           specs.UpdateProfileRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success_for_updating_profile_details",
			profileID: 1,
			userID:    1,
			input: specs.UpdateProfileRequest{
				Profile: specs.Profile{
					Name:              "Updated Name",
					Email:             "updated.email@example.com",
					Gender:            "Male",
					Mobile:            "1234567890",
					Designation:       "Updated Designation",
					Description:       "Updated Description",
					Title:             "Updated Title",
					YearsOfExperience: 5,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
					GithubLink:        "https://github.com/updated",
					LinkedinLink:      "https://linkedin.com/in/updated",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfile", mock.Anything, 1, mock.AnythingOfType("repository.UpdateProfileRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed_because_UpdateProfile_returns_an_error",
			profileID: 100000000000000000,
			userID:    1,
			input: specs.UpdateProfileRequest{
				Profile: specs.Profile{
					Name:              "Name B",
					Email:             "emailb@example.com",
					Gender:            "Female",
					Mobile:            "0987654321",
					Designation:       "Designation B",
					Description:       "Description B",
					Title:             "Title B",
					YearsOfExperience: 10,
					PrimarySkills:     []string{"Java", "C++"},
					SecondarySkills:   []string{"HTML", "CSS"},
					GithubLink:        "https://github.com/userb",
					LinkedinLink:      "https://linkedin.com/in/userb",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfile", mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProfileRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_of_missing_profile_name",
			profileID: 1,
			userID:    1,
			input: specs.UpdateProfileRequest{
				Profile: specs.Profile{
					Name:              "",
					Email:             "email@example.com",
					Gender:            "Male",
					Mobile:            "1234567890",
					Designation:       "Designation",
					Description:       "Description",
					Title:             "Title",
					YearsOfExperience: 3,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
					GithubLink:        "https://github.com/user",
					LinkedinLink:      "https://linkedin.com/in/user",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfile", mock.Anything, 1, mock.AnythingOfType("repository.UpdateProfileRepo"), mock.Anything).Return(0, errors.New("Missing profile name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_of_invalid_profile_id",
			profileID: 0,
			userID:    1,
			input: specs.UpdateProfileRequest{
				Profile: specs.Profile{
					Name:              "Name B",
					Email:             "example@gmail.com",
					Gender:            "Male",
					Mobile:            "1234567890",
					Designation:       "Designation B",
					Description:       "Description B",
					Title:             "Title B",
					YearsOfExperience: 10,
					PrimarySkills:     []string{"Java", "C++"},
					SecondarySkills:   []string{"HTML", "CSS"},
					GithubLink:        "github.com/userb",
					LinkedinLink:      "linkedin.com/in/userb",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfile", mock.Anything, 0, mock.AnythingOfType("repository.UpdateProfileRepo"), mock.Anything).Return(0, errors.New("Invalid profile id")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			_, err := profileService.UpdateProfile(context.TODO(), test.profileID, test.userID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}

			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestProfileDeleteService(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}

	profileSvc := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		profileID       int
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success_for_delete_profile",
			profileID: 1,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("DeleteProfile", mock.Anything, 1, mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed_because_DeleteProfile_returns_ErrNoData",
			profileID: 2,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("DeleteProfile", mock.Anything, 2, mock.Anything).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_DeleteProfile_returns_an_error",
			profileID: 3,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("DeleteProfile", mock.Anything, 3, mock.Anything).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_HandleTransaction_returns_an_error",
			profileID: 5,
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("DeleteProfile", mock.Anything, 5, mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, nil).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			err := profileSvc.DeleteProfile(context.Background(), test.profileID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateSequence(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		userID          int
		input           specs.UpdateSequenceRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_to_update_sequence",
			input: specs.UpdateSequenceRequest{
				ProfileID: 1,
				Component: specs.Component{
					CompName: "Component Name",
					ComponentPriorities: map[int]int{
						1: 1,
					},
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				profileMock.On("UpdateSequence", mock.Anything, mock.AnythingOfType("repository.UpdateSequenceRequest"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed_to_update_sequence",
			input: specs.UpdateSequenceRequest{
				ProfileID: 1,
				Component: specs.Component{
					CompName: "Component Name",
					ComponentPriorities: map[int]int{
						1: 1,
					},
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				profileMock.On("UpdateSequence", mock.Anything, mock.AnythingOfType("repository.UpdateSequenceRequest"), mock.Anything).Return(0, errors.New("update sequence failed")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_to_update_sequence_because_of_missing_profile_id",
			input: specs.UpdateSequenceRequest{
				ProfileID: 0,
				Component: specs.Component{
					CompName: "Component Name",
					ComponentPriorities: map[int]int{
						1: 1,
					},
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("missing profile id")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_to_update_sequence_because_of_count_records_error",
			input: specs.UpdateSequenceRequest{
				ProfileID: 1,
				Component: specs.Component{
					CompName: "Component Name",
					ComponentPriorities: map[int]int{
						1: 1,
					},
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("count records error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			_, err := profileService.UpdateSequence(context.Background(), test.userID, test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProfileStatus(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		userID          int
		input           specs.UpdateProfileStatus
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_to_update_profile_status",
			input: specs.UpdateProfileStatus{
				ProfileStatus: specs.UpdateProfileStatusRequest{
					IsCurrentEmployee: "YES",
					IsActive:          "YES",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfileStatus", mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProfileStatusRepo"), mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed_to_update_profile_status",
			input: specs.UpdateProfileStatus{
				ProfileStatus: specs.UpdateProfileStatusRequest{
					IsCurrentEmployee: "YES",
					IsActive:          "YES",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("UpdateProfileStatus", mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProfileStatusRepo"), mock.Anything).Return(errors.New("update profile status failed")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			err := profileService.UpdateProfileStatus(context.Background(), test.userID, test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
		})
	}
}
