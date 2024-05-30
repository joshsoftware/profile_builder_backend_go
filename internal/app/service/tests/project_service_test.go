package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProject(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	var repodeps = service.RepoDeps{
		ProjectDeps: mockProjectRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           dto.CreateProjectRequest
		setup           func(projectMock *mocks.ProjectStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for project details",
			input: dto.CreateProjectRequest{
				ProfileID: 1,
				Projects: []dto.Project{
					{
						Name:             "Project X",
						Description:      "Description of Project X",
						Role:             "Developer",
						Responsibilities: "Coding, testing",
						Technologies:     "Golang, Python",
						TechWorkedOn:     "Java, C++",
						WorkingStartDate: "2024-01-01",
						WorkingEndDate:   "2024-06-01",
						Duration:         "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateProject",
			input: dto.CreateProjectRequest{
				ProfileID: 123,
				Projects: []dto.Project{
					{
						Name:             "",
						Description:      "Description of Project Y",
						Role:             "Tester",
						Responsibilities: "Testing",
						Technologies:     "Java, Selenium",
						TechWorkedOn:     "Python, C#",
						WorkingStartDate: "2024-01-01",
						WorkingEndDate:   "2024-06-01",
						Duration:         "5 months",
					},
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because empty payload",
			input: dto.CreateProjectRequest{
				ProfileID: 123,
				Projects:  []dto.Project{},
			},
			setup:           func(projectMock *mocks.ProjectStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo)

			// test service
			_, err := profileService.CreateProject(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	// Initialize mock dependencies
	mockProjectRepo := new(mocks.ProjectStorer)
	var repodeps = service.RepoDeps{
		ProjectDeps: mockProjectRepo,
	}
	projectService := service.NewServices(repodeps)

	// Define mock data
	mockProfileID := "123"
	mockResponseProject := []dto.ProjectResponse{
		{
			ProfileID:        123,
			Name:             "Project Alpha",
			Description:      "A project about something",
			Role:             "Lead Developer",
			Responsibilities: "Leading the development team",
			Technologies:     "Go, React, PostgreSQL",
			TechWorkedOn:     "Backend, Frontend",
			WorkingStartDate: "2020-01-01",
			WorkingEndDate:   "2021-01-01",
			Duration:         "1 year",
		},
	}

	// Define test cases
	tests := []struct {
		name            string
		profileID       string
		setup           func(projMock *mocks.ProjectStorer)
		isErrorExpected bool
		wantResponse    []dto.ProjectResponse
	}{
		{
			name:      "Success get project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer) {
				// Mock successful retrieval
				projMock.On("GetProjects", mock.Anything, mock.Anything).Return(mockResponseProject, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseProject,
		},
		{
			name:      "Fail get project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer) {
				// Mock retrieval failure
				projMock.On("GetProjects", mock.Anything, mock.Anything).Return([]dto.ProjectResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []dto.ProjectResponse{},
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockProjectRepo)

			// Call the method being tested
			gotResp, err := projectService.GetProject(context.Background(), test.profileID)

			// Assertions
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
		})
	}
}
