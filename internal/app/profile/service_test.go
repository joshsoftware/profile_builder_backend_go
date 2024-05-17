package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfile(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateProfileRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for user Detail",
			input: dto.CreateProfileRequest{
				Profile: dto.Profile{
					Name:              "Abhishek Jain",
					Email:             "abhishek.jain@josh.com",
					Gender:            "Male",
					Mobile:            "9595601925",
					Designation:       "Software Engineer",
					Description:       "Experienced software engineer",
					Title:             "Golang Developer",
					YearsOfExperience: 5,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
					GithubLink:        "https://github.com/abhishek",
					LinkedinLink:      "https://www.linkedin.com/in/abhsihek",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed to create profile",
			input: dto.CreateProfileRequest{
				Profile: dto.Profile{
					Name:              "Abhishek Jain",
					Email:             "abhishek.jain@josh.com",
					Gender:            "Male",
					Mobile:            "9595601925",
					Designation:       "Software Engineer",
					Description:       "Experienced software engineer",
					Title:             "Golang Developer",
					YearsOfExperience: 5,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
					GithubLink:        "https://github.com/abhishek",
					LinkedinLink:      "https://www.linkedin.com/in/abhsihek",
				},
			},
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(errors.New("profile creation failed")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// test service
			err := profileService.CreateProfile(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateEducation(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateEducationRequest
		setup           func(profileMock *mocks.ProfileStorer)
		isErrorExpected bool
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
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateEducation", mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
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
			setup: func(profileMock *mocks.ProfileStorer) {
				profileMock.On("CreateEducation", mock.Anything, mock.Anything).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			err := profileService.CreateEducation(context.Background(), test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateProjectRequest
		setup           func(projectMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for project details",
			input: dto.CreateProjectRequest{
				ProfileID: 1,
				Projects: []dto.Project{
					{
						Name:             "Project X",
						Description:      "Description of Project X",
						Role:             "Developer",
						Responsibilities: "Coding, testing",
						Technologies:     "Golang, Python",
						TechWorkedOn:     "Java, C++",
						WorkingStartDate: "2024-01-01",
						WorkingEndDate:   "2024-06-01",
						Duration:         "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProfileStorer) {
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateProject",
			input: dto.CreateProjectRequest{
				ProfileID: 123,
				Projects: []dto.Project{
					{
						Name:             "",
						Description:      "Description of Project Y",
						Role:             "Tester",
						Responsibilities: "Testing",
						Technologies:     "Java, Selenium",
						TechWorkedOn:     "Python, C#",
						WorkingStartDate: "2024-01-01",
						WorkingEndDate:   "2024-06-01",
						Duration:         "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProfileStorer) {
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// test service
			err := profileService.CreateProject(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

var mockListProfile = []dto.ListProfiles{
	{
		ID:                1,
		Name:              "Abhishek Dhondalkar",
		Email:             "abhishek.dhondalkar@gmail.com",
		YearsOfExperience: 1.0,
		PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
		IsCurrentEmployee: 1,
	},
}

func TestListProfile(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

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

			// test service
			gotResp, err := profileService.ListProfiles(context.Background())
			if !assert.Equal(t, test.wantResponse, gotResp) {
				t.Errorf("Test Failed, expected resp to be %v, but got resp %v", test.wantResponse, gotResp)
			}
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
