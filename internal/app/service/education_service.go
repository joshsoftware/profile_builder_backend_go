package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// EducationService represents a set of methods for accessing the education
type EducationService interface {
	CreateEducation(ctx context.Context, eduDetail specs.CreateEducationRequest, ID int) (profileID int, err error)
	GetEducation(ctx context.Context, profileID int) (value []specs.EducationResponse, err error)
	UpdateEducation(ctx context.Context, profileID string, eduID string, req specs.UpdateEducationRequest) (id int, err error)
}

// CreateEducation : Service layer function adds education details to a user profile.
func (eduSvc *service) CreateEducation(ctx context.Context, eduDetail specs.CreateEducationRequest, ID int) (profileID int, err error) {
	if len(eduDetail.Educations) == 0 {
		zap.S().Error("educations payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.EducationRepo
	for i := 0; i < len(eduDetail.Educations); i++ {
		var val repository.EducationRepo
		val.ProfileID = ID
		val.Degree = eduDetail.Educations[i].Degree
		val.UniversityName = eduDetail.Educations[i].UniversityName
		val.Place = eduDetail.Educations[i].Place
		val.PercentageOrCgpa = eduDetail.Educations[i].PercentageOrCgpa
		val.PassingYear = eduDetail.Educations[i].PassingYear
		val.CreatedAt = today
		val.UpdatedAt = today
		//TODO by context
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = eduSvc.EducationRepo.CreateEducation(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create Education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("education(s) created with profile id : ", ID)

	return ID, nil
}

// GetEducation in the service layer retrieves a education of specific profile.
func (eduSvc *service) GetEducation(ctx context.Context, id int) (value []specs.EducationResponse, err error) {

	value, err = eduSvc.EducationRepo.GetEducation(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get education : ", err, " for profile id : ", id)
		return []specs.EducationResponse{}, err
	}
	return value, nil
}

// UpdateEducation in the service layer update a education of specific profile.
func (eduSvc *service) UpdateEducation(ctx context.Context, profileID string, eduID string, req specs.UpdateEducationRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get education params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	var value repository.UpdateEducationRepo
	value.Degree = req.Education.Degree
	value.UniversityName = req.Education.UniversityName
	value.Place = req.Education.Place
	value.PercentageOrCgpa = req.Education.PercentageOrCgpa
	value.PassingYear = req.Education.PassingYear
	value.UpdatedAt = today
	//TODO by context
	value.UpdatedByID = 1

	id, err = eduSvc.EducationRepo.UpdateEducation(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("education(s) update with profile id : ", id)

	return id, nil
}
