package post_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	userRegisterHandler := post.CreateProfileHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for user Detail",
			input: `{ "profile" : {
				"name": "Abhishek Jain",
				"email": "abhishek.jain@gmail.com",
				"gender": "Male",
				"mobile": "9595601925",
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
				mockSvc.On("CreateProfileHandler", mock.Anything).Return(nil).Once()
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
				"mobile": "9595601925",
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
				"email": "abhishek.jain@gmail.com",
				"gender": "Male",
				"mobile": "9595601925",
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
				"name": "Abhishek Jain",
				"email": "",
				"gender": "Male",
				"mobile": "9595601925",
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
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/profiles", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userRegisterHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
