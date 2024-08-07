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

func TestCreateExperienceHandler(t *testing.T) {
	ctx := context.Background()
	profileSvc := new(mocks.Service)
	createExperienceHandler := handler.CreateExperienceHandler(ctx, profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
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
				mockSvc.On("CreateExperience", mock.Anything, mock.AnythingOfType("specs.CreateExperienceRequest"), 1).Return(1, nil).Once()
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/experiences", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestListExperienceHandler(t *testing.T) {
	expSvc := mocks.NewService(t)
	getExperienceHandler := handler.ListExperienceHandler(context.Background(), expSvc)

	tests := []struct {
		name               string
		queryParams        string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success_for_getting_experiences",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetExperience", mock.Anything, 1).Return([]specs.ExperienceResponse{
					{
						ProfileID:   1,
						Designation: "Software Engineer",
						CompanyName: "Tech Corp",
						FromDate:    "2018-01-01",
						ToDate:      "2020-01-01",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail_as_error_in_getexperience",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetExperience", mock.Anything, 2).Return([]specs.ExperienceResponse{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup the mock service
			test.setup(expSvc)
			req, err := http.NewRequest("GET", "profiles/"+test.queryParams+"/experiences", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": test.queryParams})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
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
				mockSvc.On("UpdateExperience", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateExperienceRequest")).Return(1, nil).Once()
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
				mockSvc.On("UpdateExperience", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateExperienceRequest")).Return(0, errors.New("Service Error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(expSvc)

			req, err := http.NewRequest("PUT", "/profiles/1/experience/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
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
			expectedResponse:   "No data found for deletion",
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
