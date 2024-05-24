package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type projectService struct {
	ProjectRepo repository.ProjectStorer
}

type ProjectService interface {
	CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error)
}

func NewProjectService(db *pgx.Conn) ProjectService {
	projectRepo := repository.NewProjectRepo(db)

	return &projectService{
		ProjectRepo: projectRepo,
	}
}

// CreateProject : Service layer function adds project details to a user profile.
func (profileSvc *projectService) CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error) {
	if len(projDetail.Projects) == 0 {
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

	err = profileSvc.ProjectRepo.CreateProject(ctx, value)
	if err != nil {
		return 0, err
	}

	return projDetail.ProfileID, nil
}