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
	TestProjectID = 1
)

func TestCreateProjectHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	createProjectHandler := handler.CreateProjectHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
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
				mockSvc.On("CreateProject", mock.Anything, mock.AnythingOfType("specs.CreateProjectRequest"), TestProfileID, TestUserID).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Project(s) added successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : name "}`,
		},
		{
			name: "Fail_for_missing_description_field",
			input: `{
				"projects":[{
					"name": "Least Square",
					"description": "",
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
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : description "}`,
		},
		{
			name: "Failed_because_of_service_layer_error",
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
				mockSvc.On("CreateProject", mock.Anything, mock.AnythingOfType("specs.CreateProjectRequest"), TestProfileID, TestUserID).Return(0, errors.New("failed to create project")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to create project"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req := httptest.NewRequest("POST", "/profiles/projects", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProjectHandler)
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

func TestListProjectHandler(t *testing.T) {
	projSvc := mocks.NewService(t)
	getProjectHandler := handler.ListProjectHandler(context.Background(), projSvc)

	tests := []struct {
		name               string
		profileID          string
		queryParams        string
		mockSetup          func(mock *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success_for_getting_projects",
			profileID:   "1",
			queryParams: "",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListProjectsFilter{
					ProjectsIDs: nil,
					Names:       nil,
				}
				projectsResp := []specs.ProjectResponse{
					{
						ID:               1,
						ProfileID:        1,
						Name:             "Project Alpha",
						Description:      "Some Description",
						Role:             "Lead Developer",
						Responsibilities: "Developing the core features",
						Technologies:     []string{"Go", "React"},
						TechWorkedOn:     []string{"Python, C#"},
						WorkingStartDate: "2020-01-01",
						WorkingEndDate:   "2021-01-01",
						Duration:         "6 months",
					},
				}
				mockSvc.On("ListProjects", mock.Anything, 1, listFilter).Return(projectsResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"projects":[{"id":1,"profile_id":1,"name":"Project Alpha","description":"Some Description","role":"Lead Developer","responsibilities":"Developing the core features","technologies":["Go","React"],"tech_worked_on":["Python, C#"],"working_start_date":"2020-01-01","working_end_date":"2021-01-01","duration":"6 months"}]}}`,
		},
		{
			name:        "Success_for_getting_projects_with_filters",
			profileID:   "1",
			queryParams: "?projects_ids=1&names=Project%20Alpha",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListProjectsFilter{
					ProjectsIDs: []int{1},
					Names:       []string{"Project Alpha"},
				}
				projectsResp := []specs.ProjectResponse{
					{
						ID:               1,
						ProfileID:        1,
						Name:             "Project Alpha",
						Description:      "Some Description",
						Role:             "Lead Developer",
						Responsibilities: "Developing the core features",
						Technologies:     []string{"Go", "React"},
						TechWorkedOn:     []string{"Python, C#"},
						WorkingStartDate: "2020-01-01",
						WorkingEndDate:   "2021-01-01",
						Duration:         "6 months",
					},
				}
				mockSvc.On("ListProjects", mock.Anything, 1, listFilter).Return(projectsResp, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"projects":[{"id":1,"profile_id":1,"name":"Project Alpha","description":"Some Description","role":"Lead Developer","responsibilities":"Developing the core features","technologies":["Go","React"],"tech_worked_on":["Python, C#"],"working_start_date":"2020-01-01","working_end_date":"2021-01-01","duration":"6 months"}]}}`,
		},
		{
			name:      "Empty_Response_From_Service",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListProjectsFilter{
					ProjectsIDs: nil,
					Names:       nil,
				}
				mockSvc.On("ListProjects", mock.Anything, 1, listFilter).Return([]specs.ProjectResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"projects":[]}}`,
		},
		{
			name:        "Fail_as_error_in_getprojects",
			profileID:   "1",
			queryParams: "?projects_ids=2",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProjects", mock.Anything, 1, specs.ListProjectsFilter{ProjectsIDs: []int{2}}).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:      "Service_returns_Error",
			profileID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				listFilter := specs.ListProjectsFilter{
					ProjectsIDs: nil,
					Names:       nil,
				}
				mockSvc.On("ListProjects", mock.Anything, 1, listFilter).Return(nil, errors.New("service error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup(projSvc)
			req := httptest.NewRequest("GET", "/profiles/"+test.profileID+"/projects"+test.queryParams, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": test.profileID})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(test.expectedResponse) {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
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
		expectedResponse   string
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
				mockSvc.On("UpdateProject", context.Background(), TestProfileID, TestProjectID, TestUserID, mock.AnythingOfType("specs.UpdateProjectRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Project updated successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
		},
		{
			name: "Fail_for_missing_name_field",
			input: `{
				"project": {
					"name": "",
					"description": "Updated Description",
					"role": "Software Developer",
					"responsibilities": "Develop a full stack app",
					"technologies": ["Ruby, Rails, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Years"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : name"}`,
		},
		{
			name: "Fail_for_missing_description_field",
			input: `{
				"project": {
					"name": "Updated Project",
					"description": "",
					"role": "Software Developer",
					"responsibilities": "Develop a full stack app",
					"technologies": ["Ruby, Rails, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Years"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : description"}`,
		},
		{
			name: "Failed_because_of_service_layer_error",
			input: `{
				"project": {
					"name": "Updated Project",
					"description": "Updated Description",
					"role": "Software Developer",
					"responsibilities": "Develop a full stack app",
					"technologies": ["Ruby, Rails, MongoDB, AWS"],
					"tech_worked_on": ["Django, AWS"],
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Years"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateProject", context.Background(), TestProfileID, TestProjectID, TestUserID, mock.AnythingOfType("specs.UpdateProjectRequest")).Return(0, errors.New("failed to update project")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to update project"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(projSvc)
			defer projSvc.AssertExpectations(t)
			req := httptest.NewRequest("PUT", "/profiles/1/projects/1", bytes.NewBuffer([]byte(test.input)))

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateProjectHandler)
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
			expectedResponse:   `{"data":{"message":"Resource not found for the given request ID"}}`,
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
