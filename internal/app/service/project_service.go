package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ProjectService represents a set of methods for accessing the projects.
type ProjectService interface {
	CreateProject(ctx context.Context, projDetail specs.CreateProjectRequest, profileID int, userID int) (ID int, err error)
	ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter) (values []specs.ProjectResponse, err error)
	UpdateProject(ctx context.Context, profileID int, projID int, userID int, req specs.UpdateProjectRequest) (ID int, err error)
}

// CreateProject : Service layer function adds project details to a user profile.
func (projSvc *service) CreateProject(ctx context.Context, projDetail specs.CreateProjectRequest, profileID int, userID int) (ID int, err error) {
	tx, _ := projSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := projSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	if len(projDetail.Projects) == 0 {
		zap.S().Error("projects payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ProjectRepo
	for _, project := range projDetail.Projects {
		val := repository.ProjectRepo{
			ProfileID:        profileID,
			Name:             project.Name,
			Description:      project.Description,
			Role:             project.Role,
			Responsibilities: project.Responsibilities,
			Technologies:     project.Technologies,
			TechWorkedOn:     project.TechWorkedOn,
			WorkingStartDate: project.WorkingStartDate,
			WorkingEndDate:   project.WorkingEndDate,
			Duration:         project.Duration,
			CreatedAt:        today,
			UpdatedAt:        today,
			CreatedByID:      userID,
			UpdatedByID:      userID,
		}

		value = append(value, val)
	}

	err = projSvc.ProjectRepo.CreateProject(ctx, value, tx)
	if err != nil {
		zap.S().Error("Unable to create project : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("project(s) created with profile id : ", profileID)

	return profileID, nil
}

// ListProjects in the service layer retrieves a projects of specific profile.
func (projSvc *service) ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter) (values []specs.ProjectResponse, err error) {
	tx, _ := projSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := projSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	values, err = projSvc.ProjectRepo.ListProjects(ctx, profileID, filter, tx)
	if err != nil {
		zap.S().Error("Unable to get projects : ", err, " for profile id : ", profileID)
		return []specs.ProjectResponse{}, err
	}
	return values, nil
}

// UpdateProject in the service layer update a projects of specific profile.
func (projSvc *service) UpdateProject(ctx context.Context, profileID int, projID int, userID int, req specs.UpdateProjectRequest) (ID int, err error) {
	tx, _ := projSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := projSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

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
	value.UpdatedByID = userID

	profileID, err = projSvc.ProjectRepo.UpdateProject(ctx, profileID, projID, value, tx)
	if err != nil {
		zap.S().Error("Unable to update project : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("project(s) update with profile id : ", profileID)

	return profileID, nil
}
