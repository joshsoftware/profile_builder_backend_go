package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
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
				mockSvc.On("CreateProject", context.Background(), mock.AnythingOfType("dto.CreateProjectRequest")).Return(1, nil)
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
