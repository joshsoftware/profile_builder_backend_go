package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ExperienceService represents a set of methods for accessing the experiences
type ExperienceService interface {
	CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, id int) (profileID int, err error)
	GetExperience(ctx context.Context, id int) (values []specs.ExperienceResponse, err error)
	UpdateExperience(ctx context.Context, profileID string, eduID string, req specs.UpdateExperienceRequest) (id int, err error)
}

// CreateExperience : Service layer function adds experiences details to a user profile.
func (expSvc *service) CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, id int) (profileID int, err error) {
	if len(expDetail.Experiences) == 0 {
		zap.S().Error("experiences payload can't be empty")
		return 0, errors.ErrEmptyPayload
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
		//TODO by context
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = expSvc.ExperienceRepo.CreateExperience(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create experiences : ", err, " for profile id : ", id)
		return 0, err
	}
	zap.S().Info("experience(s) created with profile id : ", id)

	return id, nil
}

// GetExperience in the service layer retrieves a experiences of specific profile.
func (expSvc *service) GetExperience(ctx context.Context, id int) (values []specs.ExperienceResponse, err error) {

	values, err = expSvc.ExperienceRepo.GetExperiences(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get experiences : ", err, " for profile id : ", id)
		return []specs.ExperienceResponse{}, err
	}
	return values, nil
}

// UpdateExperience in the service layer update a experience of specific profile.
func (expSvc *service) UpdateExperience(ctx context.Context, profileID string, eduID string, req specs.UpdateExperienceRequest) (id int, err error) {
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
	//TODO by context
	value.UpdatedByID = 1

	id, err = expSvc.ExperienceRepo.UpdateExperience(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("experience(s) update with profile id : ", id)

	return id, nil
}
