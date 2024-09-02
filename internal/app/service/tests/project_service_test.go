package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
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
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
		ProjectDeps: mockProjectRepo,
	}
	profileService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           specs.CreateProjectRequest
		profileID       int
		userID          int
		setup           func(*mocks.ProjectStorer, *mocks.ProfileStorer)
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo"), mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo"), mock.Anything).Return(errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_empty_payload",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{},
			},
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo"), mock.Anything).Return(errors.New("invalid request body")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_countrecords_returns_an_error",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{
					{
						Name:             "Project Z",
						Description:      "Description of Project Z",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_invalid_profileID",
			input: specs.CreateProjectRequest{
				Projects: []specs.Project{
					{
						Name:             "Project Z",
						Description:      "Description of Project Z",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Once()
				projectMock.On("CreateProject", mock.Anything, mock.AnythingOfType("[]repository.ProjectRepo"), mock.Anything).Return(errors.New("invalid profileID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo, mockProfileRepo)
			_, err := profileService.CreateProject(context.TODO(), test.input, 1, 1)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
			mockProjectRepo.AssertExpectations(t)
			mockProfileRepo.AssertExpectations(t)
		})
	}
}

func TestListProjects(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
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
		setup           func(*mocks.ProjectStorer, *mocks.ProfileStorer)
		isErrorExpected bool
		wantResponse    []specs.ProjectResponse
	}{
		{
			name:      "Success_get_project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projMock.On("ListProjects", mock.Anything, profileID, mock.Anything, mock.Anything).Return(mockResponseProject, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseProject,
		},
		{
			name:      "Fail_get_project",
			profileID: mockProfileID,
			setup: func(projMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projMock.On("ListProjects", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.ProjectResponse{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.ProjectResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo, mockProfileRepo)
			gotResp, err := projectService.ListProjects(context.Background(), test.profileID, mockListProjectsResp)
			assert.Equal(t, test.wantResponse, gotResp)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err)
			}
			mockProfileRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	mockProjectRepo := new(mocks.ProjectStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps: mockProfileRepo,
		ProjectDeps: mockProjectRepo,
	}
	projService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		projectID       int
		userID          int
		input           specs.UpdateProjectRequest
		setup           func(*mocks.ProjectStorer, *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success_for_updating_project_details",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("UpdateProject", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed_because_UpdateProject_returns_an_error",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateProjectRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_of_missing_project_name",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("UpdateProject", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateProjectRepo"), mock.Anything).Return(0, errors.New("Missing project name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_of_invalid_profileID_or_projectID",
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
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("UpdateProject", mock.Anything, -1, 1, mock.AnythingOfType("repository.UpdateProjectRepo"), mock.Anything).Return(0, errors.New("Invalid profileID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectRepo, mockProfileRepo)
			_, err := projService.UpdateProject(context.TODO(), test.profileID, test.projectID, test.userID, test.input)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteProjectService(t *testing.T) {
	mockProjectSvc := new(mocks.ProjectStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		ProjectDeps: mockProjectSvc,
		ProfileDeps: mockProfileRepo,
	}
	projectSvc := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		projectID       int
		profileID       int
		setup           func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:      "Success_for_delete_project",
			projectID: 1,
			profileID: 1,
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("DeleteProject", mock.Anything, 1, 1, nil).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:      "Failed_because_delete_project_returns_an_error",
			projectID: 2,
			profileID: 1,
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("DeleteProject", mock.Anything, 1, 2, nil).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:      "Failed_because_DeleteProject_returns_an_error",
			projectID: 3,
			profileID: 1,
			setup: func(projectMock *mocks.ProjectStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				projectMock.On("DeleteProject", mock.Anything, 1, 3, nil).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, nil, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProjectSvc, mockProfileRepo)
			err := projectSvc.DeleteProject(context.Background(), test.profileID, test.projectID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
			mockProjectSvc.AssertExpectations(t)
			mockProfileRepo.AssertExpectations(t)
		})
	}
}
