package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
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
			name: "Success_for_education_detail",
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
				mockSvc.On("CreateEducation", mock.Anything, mock.AnythingOfType("specs.CreateEducationRequest"), 1).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_degree_field",
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
			name: "Fail_for_missing_passing_year_field",
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

func TestListEducationHandler(t *testing.T) {
	eduSvc := mocks.NewService(t)
	getEducationHandler := handler.ListEducationHandler(context.Background(), eduSvc)

	tests := []struct {
		name               string
		queryParams        string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success_for_getting_education",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetEducation", mock.Anything, 1).Return([]specs.EducationResponse{
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
			name:        "Fail_as_error_in_geteducation",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetEducation", mock.Anything, 2).Return([]specs.EducationResponse{}, errors.New("error")).Once()
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
			name: "Success_for_updating_education_detail",
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
				mockSvc.On("UpdateEducation", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateEducationRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_degree_field",
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
			name: "Fail_for_missing_university_name_field",
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
			name: "Fail_for_service_error",
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
				mockSvc.On("UpdateEducation", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateEducationRequest")).Return(0, errors.New("Service Error")).Once()
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

func TestDeleteEducationHandler(t *testing.T) {
	educationSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		educationID        string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success_for_education_delete",
			profileID:   "1",
			educationID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteEducation", mock.Anything, 1, 1).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Education deleted successfully",
		},
		{
			name:        "No_data_found_for_deletion",
			profileID:   "1",
			educationID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteEducation", mock.Anything, 1, 2).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "No data found for deletion",
		},
		{
			name:        "Error_while_deleting_education",
			profileID:   "1",
			educationID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteEducation", mock.Anything, 1, 3).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			educationID:        "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
		{
			name:        "Fail_for_internal_error",
			profileID:   "1",
			educationID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteEducation", mock.Anything, 1, 1).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(educationSvc)
			reqPath := "/profiles/" + tt.profileID + "/educations/" + tt.educationID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID, "id": tt.educationID})
			rr := httptest.NewRecorder()

			handler := handler.DeleteEducationHandler(context.Background(), educationSvc)
			handler(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if !strings.Contains(string(body), tt.expectedResponse) {
				t.Errorf("expected response to contain %q, got %q", tt.expectedResponse, body)
			}
		})
	}

}
