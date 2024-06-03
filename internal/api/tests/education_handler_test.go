package test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

func TestCreateEducationHandler(t *testing.T) {
	userSvc := new(mocks.Service)
	createEducationHandler := handler.CreateEducationHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for education detail",
			input: `{
				"educations":[{
					"degree": "BSc in Data Science",
					"university_name": "Shivaji University",
					"place": "Kolhapur",
					"percent_or_cgpa": "90.50%",
					"passing_year": "2020"
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateEducation", mock.Anything, mock.AnythingOfType("dto.CreateEducationRequest"), mock.AnythingOfType("string")).Return(1, nil).Once()
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
			name: "Fail for missing passing_year field",
			input: `{
				"educations":[{
					"degree": "BSc in Data Science",
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
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetEducationHandler(t *testing.T) {
	eduSvc := mocks.NewService(t)
	getEducationHandler := handler.GetEducationHandler(context.Background(), eduSvc)

	tests := []struct {
		name               string
		queryParams        string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success for getting education",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetEducation", mock.Anything, "1").Return([]dto.EducationResponse{
					{
						ProfileID:        1,
						Degree:           "BSc Computer Science",
						UniversityName:   "Example University",
						Place:            "City A",
						PercentageOrCgpa: "3.8",
						PassingYear:      "2015",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail as error in GetEducation",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetEducation", mock.Anything, "2").Return([]dto.EducationResponse{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(eduSvc)
			req, err := http.NewRequest("GET", "profiles/"+test.queryParams+"/educations", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": test.queryParams})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}
		})
	}
}

func TestUpdateEducationHandler(t *testing.T) {
	eduSvc := new(mocks.Service)
	updateEducationHandler := handler.UpdateEducationHandler(context.Background(), eduSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for updating education detail",
			input: `{
				"education":{
					  "degree": "MS in CS",
					  "university_name": "Cambridge University",
					  "place": "London",
					  "percent_or_cgpa": "87.50%",
					  "passing_year": "2005"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateEducation", context.Background(), "1", "1", mock.AnythingOfType("dto.UpdateEducationRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
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
				"education": {
					"degree": "",
					"university_name": "Updated University",
					"place": "Updated Place",
					"percentage_or_cgpa": "Updated CGPA",
					"passing_year": 2024
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing university_name field",
			input: `{
				"education": {
					"degree": "Updated Degree",
					"university_name": "",
					"place": "Updated Place",
					"percentage_or_cgpa": "Updated CGPA",
					"passing_year": 2024
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for service error",
			input: `{
				"education":{
					  "degree": "MS in CS",
					  "university_name": "Cambridge University",
					  "place": "London",
					  "percent_or_cgpa": "87.50%",
					  "passing_year": "2005"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateEducation", context.Background(), "1", "1", mock.AnythingOfType("dto.UpdateEducationRequest")).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(eduSvc)

			req, err := http.NewRequest("PUT", "/profiles/1/education/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
