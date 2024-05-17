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
	profileSvc := mocks.NewService(t)
	createProfileHandler := post.CreateProfileHandler(context.Background(), profileSvc)

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
				mockSvc.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(nil).Once()
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

func TestCreateEducationHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	userRegisterHandler := post.CreateEducationHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for education Detail",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "BSc in Data Science",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateEducationHandler", mock.Anything).Return(nil).Once()
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
			name: "Fail for missing degree field",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"educations":[{
					  "degree": "BSc in Data Science",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing passing_year field",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": ""
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/profiles/educations", bytes.NewBuffer([]byte(test.input)))
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

func TestCreateProjectHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	createProjectHandler := post.CreateEducationHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for project Detail",
			input: `{
				"profile_id": 1,
				"projects":[{
					"name": "Least Square",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProjectHandler", mock.Anything).Return(nil).Once()
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
			name: "Fail for missing name field",
			input: `{
				"profile_id": 1,
				"projects":[{
					"name": "",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"projects":[{
					"name": "Least Square",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/profiles/projects", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
