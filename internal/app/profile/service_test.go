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
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(1, nil).Once()
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
				profileMock.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(0, errors.New("profile creation failed")).Once()
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
		{
			name: "Failed because of empty payload",
			input: dto.CreateEducationRequest{
				ProfileID:  456,
				Educations: []dto.Education{},
			},
			setup:           func(profileMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			_, err := profileService.CreateEducation(context.Background(), test.input)
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
		{
			name: "Failed because empty payload",
			input: dto.CreateProjectRequest{
				ProfileID: 123,
				Projects:  []dto.Project{},
			},
			setup:           func(projectMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// test service
			_, err := profileService.CreateProject(context.TODO(), test.input)

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

func TestCreateExperience(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateExperienceRequest
		setup           func(experienceMock *mocks.ProfileStorer)
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
			setup: func(experienceMock *mocks.ProfileStorer) {
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
			setup: func(experienceMock *mocks.ProfileStorer) {
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
			setup: func(experienceMock *mocks.ProfileStorer) {
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
			setup:           func(experienceMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// Test the service
			_, err := profileService.CreateExperience(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateCertificate(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateCertificateRequest
		setup           func(certificateMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for certificate details",
			input: dto.CreateCertificateRequest{
				ProfileID: 1,
				Certificates: []dto.Certificate{
					{
						Name:             "Full Stack GO Certificate",
						OrganizationName: "Josh Software",
						Description:      "Description about certificate",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.ProfileStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateCertificate returns an error",
			input: dto.CreateCertificateRequest{
				ProfileID: 10000,
				Certificates: []dto.Certificate{
					{
						Name:             "Full Stack Data Science Certificate",
						OrganizationName: "Balambika Team",
						Description:      "Description",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.ProfileStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of missing certificate name",
			input: dto.CreateCertificateRequest{
				ProfileID: 1,
				Certificates: []dto.Certificate{
					{
						Name:             "",
						OrganizationName: "Org C",
						Description:      "Description C",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.ProfileStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(errors.New("Missing certificate name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of empty paylaod",
			input: dto.CreateCertificateRequest{
				ProfileID:    1,
				Certificates: []dto.Certificate{},
			},
			setup:           func(certificateMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			// Test the service
			_, err := profileService.CreateCertificate(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateAchievement(t *testing.T) {
	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)

	tests := []struct {
		name            string
		input           dto.CreateAchievementRequest
		setup           func(achievementMock *mocks.ProfileStorer)
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
			setup: func(achievementMock *mocks.ProfileStorer) {
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
			setup: func(achievementMock *mocks.ProfileStorer) {
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
			setup: func(achievementMock *mocks.ProfileStorer) {
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
			setup:           func(achievementMock *mocks.ProfileStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)

			_, err := profileService.CreateAchievement(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
