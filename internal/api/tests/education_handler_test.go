package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/mock"
)

var (
	TestProfileID   = 1
	TestUserID      = 1
	TestEducationID = 1
)

func TestCreateEducationHandler(t *testing.T) {
	userSvc := new(mocks.Service)
	createEducationHandler := handler.CreateEducationHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
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
				mockSvc.On("CreateEducation", mock.Anything, mock.AnythingOfType("specs.CreateEducationRequest"), TestProfileID, TestUserID).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Education(s) added successfully","profile_id":1}}`,
		},
		{
			name: "Success_for_education_detail_with_empty_passing_year",
			input: `{
				"educations":[{
					"degree": "BSc in Data Science",
					"university_name": "Shivaji University",
					"place": "Kolhapur",
					"percent_or_cgpa": "90.50%",
					"passing_year": ""
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateEducation", mock.Anything, mock.AnythingOfType("specs.CreateEducationRequest"), TestProfileID, TestUserID).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Education(s) added successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : degree "}`,
		},
		{
			name: "Fail_because_of_service_layer_error",
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
				mockSvc.On("CreateEducation", mock.Anything, mock.AnythingOfType("specs.CreateEducationRequest"), TestProfileID, TestUserID).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"Service Error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req := httptest.NewRequest("POST", "/profiles/educations", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})
			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.expectedResponse {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
			}

		})
	}
}

func TestListEducationHandler(t *testing.T) {
	eduSvc := mocks.NewService(t)
	getEducationHandler := handler.ListEducationHandler(context.Background(), eduSvc)

	tests := []struct {
		name               string
		profileID          string
		queryParams        string
		mockSetup          func(mock *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success_for_getting_education",
			profileID:   "1",
			queryParams: "",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListEducationsFilter{
					EduationsIDs: nil,
					Names:        nil,
				}
				eduResp := []specs.EducationResponse{
					{
						ID:               1,
						ProfileID:        1,
						Degree:           "Bachelor of Science",
						UniversityName:   "Example University",
						Place:            "City A",
						PercentageOrCgpa: "3.8",
						PassingYear:      "2015",
					},
				}
				mockSvc.On("ListEducations", mock.Anything, TestProfileID, listFilter).Return(eduResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"educations":[{"id":1,"profile_id":1,"degree":"Bachelor of Science","university_name":"Example University","place":"City A","percent_or_cgpa":"3.8","passing_year":"2015"}]}}`,
		},
		{
			name:        "Success_for_getting_education_with_filters",
			profileID:   "1",
			queryParams: "?educations_ids=1&names=Bachelor%20of%20Science",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListEducationsFilter{
					EduationsIDs: []int{1},
					Names:        []string{"Bachelor of Science"},
				}
				eduResp := []specs.EducationResponse{
					{
						ID:               1,
						ProfileID:        1,
						Degree:           "Bachelor of Science",
						UniversityName:   "Example University",
						Place:            "City A",
						PercentageOrCgpa: "3.8",
						PassingYear:      "2015",
					},
				}
				mockSvc.On("ListEducations", mock.Anything, 1, listFilter).Return(eduResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"educations":[{"id":1,"profile_id":1,"degree":"Bachelor of Science","university_name":"Example University","place":"City A","percent_or_cgpa":"3.8","passing_year":"2015"}]}}`,
		},
		{
			name:      "Empty_Response_From_Service",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListEducationsFilter{
					EduationsIDs: nil,
					Names:        nil,
				}
				mockSvc.On("ListEducations", mock.Anything, 1, listFilter).Return([]specs.EducationResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"educations":[]}}`,
		},
		{
			name:        "Fail_as_error_in_geteducation",
			profileID:   "1",
			queryParams: "?educations_ids=2",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListEducations", mock.Anything, 1, specs.ListEducationsFilter{EduationsIDs: []int{2}}).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:      "Service_returns_Error",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListEducationsFilter{
					EduationsIDs: nil,
					Names:        nil,
				}
				mockSvc.On("ListEducations", mock.Anything, 1, listFilter).Return(nil, errors.New("service error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   ` {"error_code":502,"error_message":"failed to fetch data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup(eduSvc)
			req := httptest.NewRequest("GET", "/profiles/"+test.profileID+"/educations"+test.queryParams, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(test.expectedResponse) {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
			}
			eduSvc.AssertExpectations(t)
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
		expectedResponse   string
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
				mockSvc.On("UpdateEducation", context.Background(), TestProfileID, TestEducationID, TestUserID, mock.AnythingOfType("specs.UpdateEducationRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Education updated successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
				mockSvc.On("UpdateEducation", context.Background(), TestProfileID, TestEducationID, TestUserID, mock.AnythingOfType("specs.UpdateEducationRequest")).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"Service Error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(eduSvc)

			req := httptest.NewRequest("PUT", "/profiles/1/education/1", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(updateEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.expectedResponse {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
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
			expectedResponse:   `{"data":{"message":"Resource not found for the given request ID"}}`,
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
			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

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

			fmt.Println("Response Body: ", string(body))
		})
	}

}
