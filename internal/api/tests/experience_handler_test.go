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
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/mock"
)

var (
	TestExperienceID = 1
)

func TestCreateExperienceHandler(t *testing.T) {
	ctx := context.Background()
	profileSvc := new(mocks.Service)
	createExperienceHandler := handler.CreateExperienceHandler(ctx, profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success_for_experience_detail",
			input: `{
				"experiences":[{
					"designation": "Associate Data Scientist",
					"company_name": "Josh Software Pvt.Ltd.",
					"from_date": "Jan-2023",
					"to_date": "July-2024"
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateExperience", mock.Anything, mock.AnythingOfType("specs.CreateExperienceRequest"), TestProfileID, TestUserID).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Experience(s) added successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
		},
		{
			name: "Fail_for_missing_designation_field",
			input: `{
				"experiences": [{
					"designation": "",
					"company_name": "ABC Corp",
					"from_date": "2023-01-01",
					"to_date": "2024-01-01"
				}]
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : designation "}`,
		},
		{
			name: "Fail_because_of_service_layer_error",
			input: `{
				"experiences": [{	
					"designation": "Software Engineer",	
					"company_name": "ABC Corp",
					"from_date": "2023-01-01",
					"to_date": "2024-01-01"
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateExperience", mock.Anything, mock.AnythingOfType("specs.CreateExperienceRequest"), TestProfileID, TestUserID).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"Service Error"}`,
		},
		{
			name: "Fail for invalid profile ID",
			input: `{
				"experiences": [{	
					"designation": "Software Engineer",	
					"company_name": "ABC Corp",
					"from_date": "2023-01-01",
					"to_date": "2024-01-01"
				}]
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid request data"}`,
		},
		{
			name: "Fail for missing UserID in context",
			input: `{
				"experiences": [{	
					"designation": "Software Engineer",	
					"company_name": "ABC Corp",
					"from_date": "2023-01-01",
					"to_date": "2024-01-01"
				}]
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid user id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)
			defer profileSvc.AssertExpectations(t)
			req := httptest.NewRequest("POST", "/profiles/experiences", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			if test.name == "Fail for invalid profile ID" {
				req = mux.SetURLVars(req, map[string]string{"profile_id": "invalid"})
			}

			if test.name == "Fail for missing UserID in context" {
				ctx := context.WithValue(req.Context(), constants.UserIDKey, 1)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.expectedResponse {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
			}
		})
	}
}

func TestListExperienceHandler(t *testing.T) {
	expSvc := mocks.NewService(t)
	getExperienceHandler := handler.ListExperienceHandler(context.Background(), expSvc)

	tests := []struct {
		name               string
		profileID          string
		queryParams        string
		mockSetup          func(mock *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success_for_getting_experience",
			profileID:   "1",
			queryParams: "",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListExperiencesFilter{
					ExperiencesIDs: nil,
					Names:          nil,
				}
				expResp := []specs.ExperienceResponse{
					{
						ID:          1,
						ProfileID:   1,
						Designation: "Software Engineer",
						CompanyName: "Example Corp",
						FromDate:    "2016",
						ToDate:      "2019",
					},
				}
				mockSvc.On("ListExperiences", mock.Anything, 1, listFilter).Return(expResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"experiences":[{"id":1,"profile_id":1,"designation":"Software Engineer","company_name":"Example Corp","from_date":"2016","to_date":"2019"}]}}`,
		},
		{
			name:        "Success_for_getting_experience_with_filters",
			profileID:   "1",
			queryParams: "?experiences_ids=1&names=Example%20Corp",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListExperiencesFilter{
					ExperiencesIDs: []int{1},
					Names:          []string{"Example Corp"},
				}
				expResp := []specs.ExperienceResponse{
					{
						ID:          1,
						ProfileID:   1,
						Designation: "Software Engineer",
						CompanyName: "Example Corp",
						FromDate:    "2016",
						ToDate:      "2019",
					},
				}
				mockSvc.On("ListExperiences", mock.Anything, 1, listFilter).Return(expResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"experiences":[{"id":1,"profile_id":1,"designation":"Software Engineer","company_name":"Example Corp","from_date":"2016","to_date":"2019"}]}}`,
		},
		{
			name:      "Empty_Response_From_Service",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListExperiencesFilter{
					ExperiencesIDs: nil,
					Names:          nil,
				}
				mockSvc.On("ListExperiences", mock.Anything, 1, listFilter).Return([]specs.ExperienceResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"experiences":[]}}`,
		},
		{
			name:        "Fail_as_error_in_getexperience",
			profileID:   "1",
			queryParams: "?experiences_ids=2",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListExperiences", mock.Anything, 1, specs.ListExperiencesFilter{ExperiencesIDs: []int{2}}).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:      "Service_returns_Error",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListExperiencesFilter{
					ExperiencesIDs: nil,
					Names:          nil,
				}
				mockSvc.On("ListExperiences", mock.Anything, 1, listFilter).Return(nil, errors.New("service error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:               "invalid_profile_id",
			mockSetup:          func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid request data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup(expSvc)
			req := httptest.NewRequest("GET", "/profiles/"+test.profileID+"/experiences"+test.queryParams, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": test.profileID})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			if test.name == "invalid_profile_id" {
				req = mux.SetURLVars(req, map[string]string{})
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(test.expectedResponse) {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
			}

			expSvc.AssertExpectations(t)
		})
	}
}

func TestUpdateExperienceHandler(t *testing.T) {
	expSvc := new(mocks.Service)
	updateExperienceHandler := handler.UpdateExperienceHandler(context.Background(), expSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success_for_updating_experience_detail",
			input: `{
				"experience": {
					"designation": "Updated Designation",
					"company_name": "Updated Company",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateExperience", context.Background(), TestProfileID, TestUserID, TestExperienceID, mock.AnythingOfType("specs.UpdateExperienceRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Experience updated successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
		},
		{
			name: "Fail_for_missing_designation_field",
			input: `{
				"experience": {
					"designation": "",
					"company_name": "Updated Company",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : designation "}`,
		},
		{
			name: "Fail_for_missing_company_name_field",
			input: `{
				"experience": {
					"designation": "Updated Designation",
					"company_name": "",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : company name "}`,
		},
		{
			name: "Fail_for_service_error",
			input: `{
				"experience": {
					"designation": "Updated Designation",
					"company_name": "Updated Company",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateExperience", context.Background(), TestProfileID, TestUserID, TestExperienceID, mock.AnythingOfType("specs.UpdateExperienceRequest")).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"Service Error"}`,
		},
		{
			name: "Fail_for_invalid_profile_id",
			input: `{
				"experience": {
					"designation": "Updated Designation",
					"company_name": "Updated Company",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid request data"}`,
		},
		{
			name: "Fail_for_missing_user_id_in_context",
			input: `{
				"experience": {
					"designation": "Updated Designation",
					"company_name": "Updated Company",
					"from_date": "2022-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid user id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(expSvc)
			defer expSvc.AssertExpectations(t)
			req := httptest.NewRequest("PUT", "/profiles/1/experience/1", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			if test.name == "Fail_for_invalid_profile_id" {
				req = mux.SetURLVars(req, map[string]string{"profile_id": "invalid", "id": "1"})
			}

			if test.name == "Fail_for_missing_user_id_in_context" {
				ctx := context.WithValue(req.Context(), constants.UserIDKey, 1)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateExperienceHandler)
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

func TestDeleteExperienceHandler(t *testing.T) {
	expSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		experienceID       string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:         "Success_for_experience_delete",
			profileID:    "1",
			experienceID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteExperience", mock.Anything, 1, 1).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Experience deleted successfully",
		},
		{
			name:         "No_data_found_for_deletion",
			profileID:    "1",
			experienceID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteExperience", mock.Anything, 1, 2).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Resource not found for the given request ID"}}`,
		},
		{
			name:         "Error_while_deleting_education",
			profileID:    "1",
			experienceID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteExperience", mock.Anything, 1, 3).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			experienceID:       "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
		{
			name:         "Fail_for_internal_error",
			profileID:    "1",
			experienceID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteExperience", mock.Anything, 1, 1).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(expSvc)
			reqPath := "/profiles/" + tt.profileID + "/experiences/" + tt.experienceID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID, "id": tt.experienceID})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := handler.DeleteExperienceHandler(context.Background(), expSvc)
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
