package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
				mockSvc.On("CreateEducation", context.Background(), mock.AnythingOfType("dto.CreateEducationRequest")).Return(1, nil)
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
		// {
		// 	name:        "Fail as error in GetEducation",
		// 	queryParams: "2",
		// 	setup: func(mockSvc *mocks.Service) {
		// 		mockSvc.On("GetEducation", mock.Anything, "2").Return(dto.ResponseEducation{}, errors.New("error")).Once()
		// 	},
		// 	expectedStatusCode: http.StatusBadGateway,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(eduSvc)

			req, err := http.NewRequest("GET", "profiles/"+test.queryParams+"/educations", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}
		})
	}
}
