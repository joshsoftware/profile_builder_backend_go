package test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	createProfileHandler := handler.CreateProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for user Detail",
			input: `{ "profile" : {
                "name": "Example User",
                "email": "example.user@gmail.com",
                "gender": "Male",
                "mobile": "8888999955",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing first_name field",
			input: `
                "mobile": "9999888855",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{ "profile" : {
                "name": "",
                "email": "example.user@gmail.com",
                "gender": "Male",
                "mobile": "9999888855",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing email field",
			input: `{ "profile" : {
                "name": "Example User",
                "email": "",
                "gender": "Male",
                "mobile": "9955995566",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProfileHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
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

func TestGetProfileListHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	getProfileListHandler := handler.ProfileListHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for listing profiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(mockListProfile, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail as error in ListProfiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/list_profiles", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProfileListHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	getProfileHandler := handler.GetProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		queryParams        string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success for getting profile",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProfile", mock.Anything, "1").Return(dto.ResponseProfile{
					ProfileID:         1,
					Name:              "Exaple User",
					Email:             "example@example.com",
					Gender:            "Male",
					Mobile:            "1234567890",
					Designation:       "Software Engineer",
					Description:       "Experienced software engineer",
					Title:             "Full Stack Developer",
					YearsOfExperience: 5.5,
					PrimarySkills:     []string{"Go", "JavaScript"},
					SecondarySkills:   []string{"Python", "Java"},
					GithubLink:        "https://github.com/demo",
					LinkedinLink:      "https://linkedin.com/in/demo",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail as error in GetProfile",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProfile", mock.Anything, "2").Return(dto.ResponseProfile{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/profile?id="+test.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProfileHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}
		})
	}
}
