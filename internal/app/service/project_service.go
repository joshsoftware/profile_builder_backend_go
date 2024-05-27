package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ProjectService represents a set of methods for accessing the projects.
type ProjectService interface {
	CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error)
	GetProject(ctx context.Context, profileID string) (values []dto.ProjectResponse, err error)
}

// CreateProject : Service layer function adds project details to a user profile.
func (projSvc *service) CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error) {
	if len(projDetail.Projects) == 0 {
		zap.S().Error("projects payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ProjectDao
	for i := 0; i < len(projDetail.Projects); i++ {
		var val repository.ProjectDao

		val.ProfileID = projDetail.ProfileID
		val.Name = projDetail.Projects[i].Name
		val.Description = projDetail.Projects[i].Description
		val.Role = projDetail.Projects[i].Role
		val.Responsibilities = projDetail.Projects[i].Responsibilities
		val.Technologies = projDetail.Projects[i].Technologies
		val.TechWorkedOn = projDetail.Projects[i].TechWorkedOn
		val.Duration = projDetail.Projects[i].Duration
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = projSvc.ProjectRepo.CreateProject(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create project : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return projDetail.ProfileID, nil
}

// GetProject in the service layer retrieves a projects of specific profile.
func (projSvc *service) GetProject(ctx context.Context, profileID string) (values []dto.ProjectResponse, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err != nil {
		zap.S().Error("error to get projects params : ", err, " for profile id : ", profileID)
		return []dto.ProjectResponse{}, err
	}

	values, err = projSvc.ProjectRepo.GetProjects(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get projects : ", err, " for profile id : ", profileID)
		return []dto.ProjectResponse{}, err
	}
	return values, nil
}
