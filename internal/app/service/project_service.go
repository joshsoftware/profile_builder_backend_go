package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ProjectService represents a set of methods for accessing the projects.
type ProjectService interface {
	CreateProject(ctx context.Context, projDetail specs.CreateProjectRequest, id int) (profileID int, err error)
	GetProject(ctx context.Context, profileID int) (values []specs.ProjectResponse, err error)
	UpdateProject(ctx context.Context, profileID string, eduID string, req specs.UpdateProjectRequest) (id int, err error)
}

// CreateProject : Service layer function adds project details to a user profile.
func (projSvc *service) CreateProject(ctx context.Context, projDetail specs.CreateProjectRequest, id int) (profileID int, err error) {
	if len(projDetail.Projects) == 0 {
		zap.S().Error("projects payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ProjectRepo
	for i := 0; i < len(projDetail.Projects); i++ {
		var val repository.ProjectRepo

		val.ProfileID = id
		val.Name = projDetail.Projects[i].Name
		val.Description = projDetail.Projects[i].Description
		val.Role = projDetail.Projects[i].Role
		val.Responsibilities = projDetail.Projects[i].Responsibilities
		val.Technologies = projDetail.Projects[i].Technologies
		val.TechWorkedOn = projDetail.Projects[i].TechWorkedOn
		val.WorkingStartDate = projDetail.Projects[i].WorkingStartDate
		val.WorkingEndDate = projDetail.Projects[i].WorkingEndDate
		val.Duration = projDetail.Projects[i].Duration
		val.CreatedAt = today
		val.UpdatedAt = today
		//TODO by context
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = projSvc.ProjectRepo.CreateProject(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create project : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("project(s) created with profile id : ", id)

	return id, nil
}

// GetProject in the service layer retrieves a projects of specific profile.
func (projSvc *service) GetProject(ctx context.Context, profileID int) (values []specs.ProjectResponse, err error) {
	values, err = projSvc.ProjectRepo.GetProjects(ctx, profileID)
	if err != nil {
		zap.S().Error("Unable to get projects : ", err, " for profile id : ", profileID)
		return []specs.ProjectResponse{}, err
	}
	return values, nil
}

// UpdateProject in the service layer update a projects of specific profile.
func (projSvc *service) UpdateProject(ctx context.Context, profileID string, eduID string, req specs.UpdateProjectRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get project params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	var value repository.UpdateProjectRepo
	value.Name = req.Project.Name
	value.Description = req.Project.Description
	value.Role = req.Project.Role
	value.Responsibilities = req.Project.Responsibilities
	value.Technologies = req.Project.Technologies
	value.TechWorkedOn = req.Project.TechWorkedOn
	value.WorkingStartDate = req.Project.WorkingStartDate
	value.WorkingEndDate = req.Project.WorkingEndDate
	value.Duration = req.Project.Duration
	value.UpdatedAt = today
	//TODO by context
	value.UpdatedByID = 1

	id, err = projSvc.ProjectRepo.UpdateProject(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update project : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("project(s) update with profile id : ", id)

	return id, nil
}
