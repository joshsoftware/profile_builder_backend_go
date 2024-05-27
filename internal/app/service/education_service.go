package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// EducationService represents a set of methods for accessing the education
type EducationService interface {
	CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error)
	GetEducation(ctx context.Context, profileID string) (value []dto.EducationResponse, err error)
	UpdateEducation(ctx context.Context, profileID string, eduID string, req dto.UpdateEducationRequest) (id int, err error)
}

// CreateEducation : Service layer function adds education details to a user profile.
func (eduSvc *service) CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error) {
	if len(eduDetail.Educations) == 0 {
		zap.S().Error("educations payload can't be empty")
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

	err = eduSvc.EducationRepo.CreateEducation(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create Education : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return eduDetail.ProfileID, nil
}

// GetEducation in the service layer retrieves a education of specific profile.
func (eduSvc *service) GetEducation(ctx context.Context, profileID string) (value []dto.EducationResponse, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err != nil {
		zap.S().Error("error to get education : ", err, " for profile id : ", profileID)
		return []dto.EducationResponse{}, err
	}

	value, err = eduSvc.EducationRepo.GetEducation(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get education : ", err, " for profile id : ", profileID)
		return []dto.EducationResponse{}, err
	}
	return value, nil
}

// UpdateEducation in the service layer update a education of specific profile.
func (eduSvc *service) UpdateEducation(ctx context.Context, profileID string, eduID string, req dto.UpdateEducationRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get education params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	id, err = eduSvc.EducationRepo.UpdateEducation(ctx, pid, id, req)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	return id, nil
}
