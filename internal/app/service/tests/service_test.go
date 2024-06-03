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

var mockListProfile = []dto.ListProfiles{
	{
		ID:                1,
		Name:              "Abhishek Dhondalkar",
		Email:             "abhishek.dhondalkar@gmail.com",
		YearsOfExperience: 1.0,
		PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
		IsCurrentEmployee: "YES",
	},
}

var mockProfile = dto.Profile{
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

var mockProfileID = "123"

var mockResponseProfile = dto.ResponseProfile{
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
	LinkedinLink:      "https://www.linkedin.com/in/abhsihek",
}

var mockListSkills = dto.ListSkills{Name: []string{"GO", "RUBY", "C", "C++", "JAVA", "PYTHON", "JAVASCRIPT"}}

func TestListProfile(t *testing.T) {
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		setup           func(userMock *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    []dto.ListProfiles
	}{
		{
			name: "Success get list of Profiles",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("ListProfiles", mock.Anything).Return(mockListProfile, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockListProfile,
		},
		{
			name: "Fail get list of Profiles",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("ListProfiles", mock.Anything).Return([]dto.ListProfiles{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []dto.ListProfiles{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// Test service
			gotResp, err := profileService.ListProfiles(context.Background())

			assert.Equal(t, test.wantResponse, gotResp)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
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
		wantResponse    dto.ListSkills
	}{
		{
			name: "Success get list of Skills",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("ListSkills", mock.Anything).Return(mockListSkills, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockListSkills,
		},
		{
			name: "Fail get list of Skills",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("ListSkills", mock.Anything).Return(dto.ListSkills{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    dto.ListSkills{},
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
		input           dto.CreateProfileRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for user Detail",
			input: dto.CreateProfileRequest{
				Profile: mockProfile,
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("ProfileRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed to create profile",
			input: dto.CreateProfileRequest{
				Profile: mockProfile,
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("ProfileRepo")).Return(0, errors.New("profile creation failed")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// test service
			_, err := profileService.CreateProfile(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
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
		profileID       string
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    dto.ResponseProfile
	}{
		{
			name:      "Success get profile",
			profileID: mockProfileID,
			setup: func(profileMock *mocks.ProfileStorer) {
				// Mock successful retrieval
				profileMock.On("GetProfile", mock.Anything, mock.Anything).Return(mockResponseProfile, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseProfile,
		},
		{
			name:      "Fail get profile",
			profileID: mockProfileID,
			setup: func(profileMock *mocks.ProfileStorer) {
				// Mock retrieval failure
				profileMock.On("GetProfile", mock.Anything, mock.Anything).Return(dto.ResponseProfile{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    dto.ResponseProfile{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockProfileRepo)

			// Call the method being tested
			gotResp, err := profileService.GetProfile(context.Background(), test.profileID)

			// Assertions
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
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
		profileID       string
		input           dto.UpdateProfileRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success for updating profile details",
			profileID: "1",
			input: dto.UpdateProfileRequest{
				Profile: dto.Profile{
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
				profileMock.On("UpdateProfile", mock.Anything, 1, mock.AnythingOfType("repository.UpdateProfileRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed because UpdateProfile returns an error",
			profileID: "100000000000000000",
			input: dto.UpdateProfileRequest{
				Profile: dto.Profile{
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
				profileMock.On("UpdateProfile", mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProfileRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of missing profile name",
			profileID: "1",
			input: dto.UpdateProfileRequest{
				Profile: dto.Profile{
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
				profileMock.On("UpdateProfile", mock.Anything, 1, mock.AnythingOfType("repository.UpdateProfileRepo")).Return(0, errors.New("Missing profile name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of invalid profileID",
			profileID: "invalid",
			input: dto.UpdateProfileRequest{
				Profile: dto.Profile{
					Name:              "Valid Name",
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
			setup:           func(profileMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			_, err := profileService.UpdateProfile(context.TODO(), test.profileID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
