package profile

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type service struct {
	Repo repository.ProfileStorer
}

type Service interface {
	CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) error
	CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) error
	CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) error
	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
}

func NewServices(repo repository.ProfileStorer) Service {
	return &service{
		Repo: repo,
	}
}

func (profileSvc *service) CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) error {

	err := profileSvc.Repo.CreateProfile(ctx, profileDetail)
	if err != nil {
		return err
	}

	return nil
}

func (profileSvc *service) CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) error {
	today := helpers.GetTodaysDate()

	var value []repository.EducationDao
	for i := 0; i < len(eduDetail.Educations); i++ {
		var val repository.EducationDao
		val.ProfileID = eduDetail.ProfileID
		val.Degree = eduDetail.Educations[i].Degree
		val.UniversityName = eduDetail.Educations[i].UniversityName
		val.Place = eduDetail.Educations[i].Place
		val.PercentageOrCgpa = eduDetail.Educations[i].PercentageOrCgpa
		val.PassingYear = eduDetail.Educations[i].PassingYear
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err := profileSvc.Repo.CreateEducation(ctx, value)
	if err != nil {
		return err
	}

	return nil
}

func (profileSvc *service) CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) error {
	today := helpers.GetTodaysDate()

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

	err := profileSvc.Repo.CreateProject(ctx, value)
	if err != nil {
		return err
	}

	return nil
}

func (profileSvc *service) ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error) {

	values, err = profileSvc.Repo.ListProfiles(ctx)
	if err != nil {
		return []dto.ListProfiles{}, err
	}

	return values, nil
}
