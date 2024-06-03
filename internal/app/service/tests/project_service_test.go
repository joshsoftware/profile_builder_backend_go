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
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateProject",
			input: dto.CreateProjectRequest{
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
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because empty payload",
			input: dto.CreateProjectRequest{
				Projects: []dto.Project{},
			},
			setup:           func(projectMock *mocks.ProjectStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo)

			// test service
			_, err := profileService.CreateProject(context.TODO(), test.input, "1")

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

func TestUpdateProject(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	var repodeps = service.RepoDeps{
		ProjectDeps: mockProjectRepo,
	}
	projService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       string
		projectID       string
		input           dto.UpdateProjectRequest
		setup           func(projectMock *mocks.ProjectStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success for updating project details",
			profileID: "1",
			projectID: "1",
			input: dto.UpdateProjectRequest{
				Project: dto.Project{
					Name:             "Updated Project Name",
					Description:      "Updated Description",
					Role:             "Updated Role",
					Responsibilities: "Updated Responsibilities",
					Technologies:     "Updated Technologies",
					TechWorkedOn:     "Updated TechWorkedOn",
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed because UpdateProject returns an error",
			profileID: "100000000000000000",
			projectID: "1",
			input: dto.UpdateProjectRequest{
				Project: dto.Project{
					Name:             "Project B",
					Description:      "Description B",
					Role:             "Role B",
					Responsibilities: "Responsibilities B",
					Technologies:     "Technologies B",
					TechWorkedOn:     "TechWorkedOn B",
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of missing project name",
			profileID: "1",
			projectID: "1",
			input: dto.UpdateProjectRequest{
				Project: dto.Project{
					Name:             "",
					Description:      "Description",
					Role:             "Role",
					Responsibilities: "Responsibilities",
					Technologies:     "Technologies",
					TechWorkedOn:     "TechWorkedOn",
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(0, errors.New("Missing project name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of invalid profileID or projectID",
			profileID: "invalid",
			projectID: "1",
			input: dto.UpdateProjectRequest{
				Project: dto.Project{
					Name:             "Valid Name",
					Description:      "Valid Description",
					Role:             "Valid Role",
					Responsibilities: "Valid Responsibilities",
					Technologies:     "Valid Technologies",
					TechWorkedOn:     "Valid TechWorkedOn",
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup:           func(projectMock *mocks.ProjectStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo)

			_, err := projService.UpdateProject(context.TODO(), test.profileID, test.projectID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
