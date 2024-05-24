package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type educationService struct {
	EducationRepo repository.EducationStorer
}

type EducationService interface {
	CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error)
}

func NewEducationService(db *pgx.Conn) EducationService {
	achievementRepo := repository.NewEducationRepo(db)

	return &educationService{
		EducationRepo: achievementRepo,
	}
}

// CreateEducation : Service layer function adds education details to a user profile.
func (profileSvc *educationService) CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error) {
	if len(eduDetail.Educations) == 0 {
		return 0, errors.ErrEmptyPayload
	}

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

	err = profileSvc.EducationRepo.CreateEducation(ctx, value)
	if err != nil {
		return 0, err
	}

	return eduDetail.ProfileID, nil
}