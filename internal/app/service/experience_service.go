package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type experienceService struct {
	ExperienceRepo repository.ExperienceStorer
}

type ExperienceService interface {
	CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest) (profileID int, err error)
}

func NewExperienceService(db *pgx.Conn) ExperienceService {
	experienceRepo := repository.NewExperienceRepo(db)

	return &experienceService{
		ExperienceRepo: experienceRepo,
	}
}

// CreateProject : Service layer function adds experiences details to a user profile.
func (profileSvc *experienceService) CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest) (profileID int, err error) {
	if len(expDetail.Experiences) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ExperienceDao
	for i := 0; i < len(expDetail.Experiences); i++ {
		var val repository.ExperienceDao

		val.ProfileID = expDetail.ProfileID
		val.Designation = expDetail.Experiences[i].Designation
		val.CompanyName = expDetail.Experiences[i].CompanyName
		val.FromDate = expDetail.Experiences[i].FromDate
		val.ToDate = expDetail.Experiences[i].ToDate
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.ExperienceRepo.CreateExperience(ctx, value)
	if err != nil {
		return 0, err
	}

	return expDetail.ProfileID, nil
}
