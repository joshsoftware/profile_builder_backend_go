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
	CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest) (profileID int, err error)
	GetExperience(ctx context.Context, profileID string) (values []dto.ExperienceResponse, err error)
}

// CreateExperience : Service layer function adds experiences details to a user profile.
func (expSvc *service) CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest) (profileID int, err error) {
	if len(expDetail.Experiences) == 0 {
		zap.S().Error("experiences payload can't be empty")
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

	err = expSvc.ExperienceRepo.CreateExperience(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create experiences : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return expDetail.ProfileID, nil
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
