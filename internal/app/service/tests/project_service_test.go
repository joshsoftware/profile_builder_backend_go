package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockListProjectsResp = specs.ListProjectsFilter{
	ProjectsIDs: []int{},
	Names:       []string{},
}

func TestCreateProject(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	var repodeps = service.RepoDeps{
		ProjectDeps: mockProjectRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           specs.CreateProjectRequest
		setup           func(projectMock *mocks.ProjectStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_for_project_details",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{
					{
						Name:             "Project X",
						Description:      "Description of Project X",
						Role:             "Developer",
						Responsibilities: "Coding, testing",
						Technologies:     []string{"Java", "Selenium"},
						TechWorkedOn:     []string{"Python, C#"},
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
			name: "Failed_because_createproject",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{
					{
						Name:             "",
						Description:      "Description of Project Y",
						Role:             "Tester",
						Responsibilities: "Testing",
						Technologies:     []string{"Java", "Selenium"},
						TechWorkedOn:     []string{"Python, C#"},
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
			name: "Failed_because_empty_payload",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{},
			},
			setup:           func(projectMock *mocks.ProjectStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo)

			// test service
			_, err := profileService.CreateProject(context.TODO(), test.input, 1, 1)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestListProjects(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	var repodeps = service.RepoDeps{
		ProjectDeps: mockProjectRepo,
	}
	projectService := service.NewServices(repodeps)

	mockResponseProject := []specs.ProjectResponse{
		{
			ProfileID:        mockProfileID,
			Name:             "Project Alpha",
			Description:      "A project about something",
			Role:             "Lead Developer",
			Responsibilities: "Leading the development team",
			Technologies:     []string{"Java", "Selenium"},
			TechWorkedOn:     []string{"Python, C#"},
			WorkingStartDate: "2020-01-01",
			WorkingEndDate:   "2021-01-01",
			Duration:         "1 year",
		},
	}

	tests := []struct {
		name            string
		profileID       int
		setup           func(projMock *mocks.ProjectStorer)
		isErrorExpected bool
		wantResponse    []specs.ProjectResponse
	}{
		{
			name:      "Success_get_project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer) {
				projMock.On("ListProjects", mock.Anything, mock.Anything, mock.Anything).Return(mockResponseProject, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseProject,
		},
		{
			name:      "Fail_get_project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer) {
				projMock.On("ListProjects", mock.Anything, mock.Anything, mock.Anything).Return([]specs.ProjectResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.ProjectResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup mock
			test.setup(mockProjectRepo)

			// Call the method being tested
			gotResp, err := projectService.ListProjects(context.Background(), test.profileID, mockListProjectsResp)

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
		profileID       int
		projectID       int
		userID          int
		input           specs.UpdateProjectRequest
		setup           func(projectMock *mocks.ProjectStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success for updating project details",
			profileID: 1,
			projectID: 1,
			userID:    1,
			input: specs.UpdateProjectRequest{
				Project: specs.Project{
					Name:             "Updated Project Name",
					Description:      "Updated Description",
					Role:             "Updated Role",
					Responsibilities: "Updated Responsibilities",
					Technologies:     []string{"Java", "Selenium"},
					TechWorkedOn:     []string{"Python, C#"},
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed because UpdateProject returns an error",
			profileID: 10000000,
			projectID: 1,
			userID:    1,
			input: specs.UpdateProjectRequest{
				Project: specs.Project{
					Name:             "Project B",
					Description:      "Description B",
					Role:             "Role B",
					Responsibilities: "Responsibilities B",
					Technologies:     []string{"Java", "Selenium"},
					TechWorkedOn:     []string{"Python, C#"},
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of missing project name",
			profileID: 1,
			projectID: 1,
			userID:    1,
			input: specs.UpdateProjectRequest{
				Project: specs.Project{
					Name:             "",
					Description:      "Description",
					Role:             "Role",
					Responsibilities: "Responsibilities",
					Technologies:     []string{"Java", "Selenium"},
					TechWorkedOn:     []string{"Python, C#"},
					WorkingStartDate: "2022-01-01",
					WorkingEndDate:   "2023-01-01",
					Duration:         "1 year",
				},
			},
			setup: func(projectMock *mocks.ProjectStorer) {
				projectMock.On("UpdateProject", mock.Anything, 1, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo")).Return(0, errors.New("Missing project name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed because of invalid profileID or projectID",
			profileID: -1,
			projectID: 1,
			userID:    1,
			input: specs.UpdateProjectRequest{
				Project: specs.Project{
					Name:             "Valid Name",
					Description:      "Valid Description",
					Role:             "Valid Role",
					Responsibilities: "Valid Responsibilities",
					Technologies:     []string{"Java", "Selenium"},
					TechWorkedOn:     []string{"Python, C#"},
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

			_, err := projService.UpdateProject(context.TODO(), test.profileID, test.projectID, test.userID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
