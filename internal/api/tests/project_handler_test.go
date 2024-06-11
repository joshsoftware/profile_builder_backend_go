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
			name: "Success for project Detail",
			input: `{
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
				mockSvc.On("CreateProject", mock.Anything, mock.AnythingOfType("dto.CreateProjectRequest"), "1").Return(1, nil).Once()
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
			name:        "Success for getting projects",
			queryParams: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProject", mock.Anything, "1").Return([]specs.ProjectResponse{
					{
						ProfileID:        1,
						Name:             "Project Alpha",
						Description:      "A sample project",
						Role:             "Lead Developer",
						Responsibilities: "Developing the core features",
						Technologies:     "Go, React",
						TechWorkedOn:     "Go, React, Docker",
						WorkingStartDate: "2020-01-01",
						WorkingEndDate:   "2021-01-01",
						Duration:         "1 year",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail as error in GetProject",
			queryParams: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProject", mock.Anything, "2").Return([]specs.ProjectResponse{}, errors.New("error")).Once()
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
			name: "Success for updating project detail",
			input: `{
				"project":{
					"name": "Profile Builder",
					"description": "A Webapp Which is Used to Build a Standard IOT Apps",
					"role": "Software Developer",
					"responsibilities": "Develop a full stack app",
					"technologies": "Ruby, Rails, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Years"
					}
				}
				`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateProject", context.Background(), "1", "1", mock.AnythingOfType("dto.UpdateProjectRequest")).Return(1, nil).Once()
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
			name: "Fail for missing name field",
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
			name: "Fail for missing description field",
			input: `{
				"project": {
					"name": "Updated Project",
					"description": ""
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		// {
		// 	name: "Fail for service error",
		// 	input: `{
		// 		"project": {
		// 			"name": "Updated Project",
		// 			"description": "Updated Description"
		// 		}
		// 	}`,
		// 	setup: func(mockSvc *mocks.Service) {
		// 		mockSvc.On("UpdateProject", context.Background(), "1", "1", mock.AnythingOfType("dto.UpdateProjectRequest")).Return(0, errors.New("Service Error")).Once()
		// 	},
		// 	expectedStatusCode: http.StatusBadGateway,
		// },
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
