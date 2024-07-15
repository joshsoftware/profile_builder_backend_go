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

func TestCreateProjectHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	createProjectHandler := handler.CreateProjectHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_project_Detail",
			input: `{
				"projects":[{
					"name": "Least Square",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": ["Python, Django, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProject", mock.Anything, mock.AnythingOfType("specs.CreateProjectRequest"), 1).Return(1, nil).Once()
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
			name: "Fail_for_missing_name_field",
			input: `{
				"projects":[{
					"name": "",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": ["Python, Django, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
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
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestListProjectHandler(t *testing.T) {
	projSvc := mocks.NewService(t)
	getProjectHandler := handler.ListProjectHandler(context.Background(), projSvc)

	tests := []struct {
		name               string
		queryParams        string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success_for_getting_projects",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProject", mock.Anything, 1).Return([]specs.ProjectResponse{
					{
						ProfileID:        1,
						Name:             "Project Alpha",
						Description:      "A sample project",
						Role:             "Lead Developer",
						Responsibilities: "Developing the core features",
						Technologies:     []string{"Java", "Selenium"},
						TechWorkedOn:     []string{"Python, C#"},
						WorkingStartDate: "2020-01-01",
						WorkingEndDate:   "2021-01-01",
						Duration:         "1 year",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail_as_error_in_getproject",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProject", mock.Anything, 2).Return([]specs.ProjectResponse{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(projSvc)

			req, err := http.NewRequest("GET", "profiles/"+test.queryParams+"/projects", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": test.queryParams})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}

			projSvc.AssertExpectations(t)
		})
	}
}

func TestUpdateProjectHandler(t *testing.T) {
	projSvc := new(mocks.Service)
	updateProjectHandler := handler.UpdateProjectHandler(context.Background(), projSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_updating_project_detail",
			input: `{
				"project":{
					"name": "Profile Builder",
					"description": "A Webapp Which is Used to Build a Standard IOT Apps",
					"role": "Software Developer",
					"responsibilities": "Develop a full stack app",
					"technologies": ["Ruby, Rails, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Years"
					}
				}
				`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateProject", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateProjectRequest")).Return(1, nil).Once()
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
			name: "Fail_for_missing_name_field",
			input: `{
				"project": {
					"name": "",
					"description": "Updated Description",
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_description_field",
			input: `{
				"project": {
					"name": "Updated Project",
					"description": ""
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(projSvc)

			req, err := http.NewRequest("PUT", "/profiles/1/projects/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestDeleteProjectHandler(t *testing.T) {
	projectSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		projectID          string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "Success_for_project_delete",
			profileID: "1",
			projectID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProject", mock.Anything, 1, 1).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Project deleted successfully",
		},
		{
			name:      "No_data_found_for_deletion",
			profileID: "1",
			projectID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProject", mock.Anything, 1, 2).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "No data found for deletion",
		},
		{
			name:      "Error_while_deleting_project",
			profileID: "1",
			projectID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProject", mock.Anything, 1, 3).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			projectID:          "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
		{
			name:      "Fail_for_internal_error",
			profileID: "1",
			projectID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProject", mock.Anything, 1, 1).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(projectSvc)
			reqPath := "/profiles/" + tt.profileID + "/projects/" + tt.projectID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID, "id": tt.projectID})
			rr := httptest.NewRecorder()

			handler := handler.DeleteProjectHandler(context.Background(), projectSvc)
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
