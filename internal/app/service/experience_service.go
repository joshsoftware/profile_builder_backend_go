package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ExperienceService represents a set of methods for accessing the experiences
type ExperienceService interface {
	CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest, ID string) (profileID int, err error)
	GetExperience(ctx context.Context, profileID string) (values []dto.ExperienceResponse, err error)
	UpdateExperience(ctx context.Context, profileID string, eduID string, req dto.UpdateExperienceRequest) (id int, err error)
}

// CreateExperience : Service layer function adds experiences details to a user profile.
func (expSvc *service) CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest, ID string) (profileID int, err error) {
	if len(expDetail.Experiences) == 0 {
		zap.S().Error("experiences payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	id, err := helpers.ConvertStringToInt(ID)
	if err != nil {
		zap.S().Error("error to get achievement params : ", err, " for profile id : ", ID)
		return 0, err
	}

	var value []repository.ExperienceRepo
	for i := 0; i < len(expDetail.Experiences); i++ {
		var val repository.ExperienceRepo

		val.ProfileID = id
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

	err = expSvc.ExperienceRepo.CreateExperience(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create experiences : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return id, nil
}

// GetExperience in the service layer retrieves a experiences of specific profile.
func (expSvc *service) GetExperience(ctx context.Context, profileID string) (values []dto.ExperienceResponse, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err != nil {
		zap.S().Error("error to get experience params : ", err, " for profile id : ", profileID)
		return []dto.ExperienceResponse{}, err
	}

	values, err = expSvc.ExperienceRepo.GetExperiences(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get experiences : ", err, " for profile id : ", profileID)
		return []dto.ExperienceResponse{}, err
	}
	return values, nil
}

// UpdateExperience in the service layer update a experience of specific profile.
func (expSvc *service) UpdateExperience(ctx context.Context, profileID string, eduID string, req dto.UpdateExperienceRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get experience params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	var value repository.UpdateExperienceRepo
	value.Designation = req.Experience.Designation
	value.CompanyName = req.Experience.CompanyName
	value.FromDate = req.Experience.FromDate
	value.ToDate = req.Experience.ToDate
	value.UpdatedAt = today
	value.UpdatedByID = 1

	id, err = expSvc.ExperienceRepo.UpdateExperience(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	return id, nil
}
