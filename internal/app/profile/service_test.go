package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfile(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateProfileRequest
		setup           func(userMock *mocks.ProfileStorer)
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
					YearsOfExperience: 5.0,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
					GithubLink:        "https://github.com/abhishek",
					LinkedinLink:      "https://www.linkedin.com/in/abhsihek",
				},
			},
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("RegisterUser", mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because RegisterUser",
			input: dto.CreateProfileRequest{
				Profile: dto.Profile{
					Name:              "",
					Email:             "@example.com",
					Gender:            "Male",
					Mobile:            "567890",
					Designation:       "Software Engineer",
					Description:       "Experienced software engineer",
					Title:             "Mr.",
					YearsOfExperience: 5.0,
					PrimarySkills:     []string{"Golang", "Python"},
					SecondarySkills:   []string{"JavaScript", "SQL"},
				},
			},
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("RegisterUser", mock.Anything).Return(errors.New("Error")).Once()
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
				ProfileId: 1,
				Educations: []dto.Education{
					{
						Degree:           "B.Tech in Computer Science",
						UniversityName:  "SPPU",
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
				ProfileId: 456,
				Educations: []dto.Education{
					{
						Degree:           "",
						UniversityName:  "Example University",
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
				ProfileId: 1,
				Projects: []dto.Project{
					{
						Name:              "Project X",
						Description:       "Description of Project X",
						Role:              "Developer",
						Responsibilities:  "Coding, testing",
						Technologies:      "Golang, Python",
						TechWorkedOn:      "Java, C++",
						WorkingStartDate:  "2024-01-01",
						WorkingEndDate:    "2024-06-01",
						Duration:          "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProfileStorer) {
				projectMock.On("CreateProject", mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateProject",
			input: dto.CreateProjectRequest{
				ProfileId: 123,
				Projects: []dto.Project{
					{
						Name:              "",
						Description:       "Description of Project Y",
						Role:              "Tester",
						Responsibilities:  "Testing",
						Technologies:      "Java, Selenium",
						TechWorkedOn:      "Python, C#",
						WorkingStartDate:  "2024-01-01",
						WorkingEndDate:    "2024-06-01",
						Duration:          "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProfileStorer) {
				projectMock.On("CreateProject", mock.Anything).Return(errors.New("Error")).Once()
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
