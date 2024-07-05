package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// EducationService represents a set of methods for accessing the education
type EducationService interface {
	CreateEducation(ctx context.Context, eduDetail specs.CreateEducationRequest, profileID int, userID int) (ID int, err error)
	ListEducations(ctx context.Context, id int, filter specs.ListEducationsFilter) (value []specs.EducationResponse, err error)
	UpdateEducation(ctx context.Context, profileID int, eduID int, userID int, req specs.UpdateEducationRequest) (ID int, err error)
}

// CreateEducation : Service layer function adds education details to a user profile.
func (eduSvc *service) CreateEducation(ctx context.Context, eduDetail specs.CreateEducationRequest, profileID int, userID int) (ID int, err error) {
	tx, _ := eduSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := eduSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	if len(eduDetail.Educations) == 0 {
		zap.S().Error("educations payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.EducationRepo
	for _, education := range eduDetail.Educations {
		val := repository.EducationRepo{
			ProfileID:        profileID,
			Degree:           education.Degree,
			UniversityName:   education.UniversityName,
			Place:            education.Place,
			PercentageOrCgpa: education.PercentageOrCgpa,
			PassingYear:      education.PassingYear,
			CreatedAt:        today,
			UpdatedAt:        today,
			CreatedByID:      userID,
			UpdatedByID:      userID,
		}

		value = append(value, val)
	}

	err = eduSvc.EducationRepo.CreateEducation(ctx, value, tx)
	if err != nil {
		zap.S().Error("Unable to create Education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("education(s) created with profile id : ", profileID)

	return profileID, nil
}

// ListEducations in the service layer retrieves a education of specific profile.
func (eduSvc *service) ListEducations(ctx context.Context, id int, filter specs.ListEducationsFilter) (value []specs.EducationResponse, err error) {
	tx, _ := eduSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := eduSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	value, err = eduSvc.EducationRepo.ListEducations(ctx, id, filter, tx)
	if err != nil {
		zap.S().Error("Unable to get education : ", err, " for profile id : ", id)
		return []specs.EducationResponse{}, err
	}
	return value, nil
}

// UpdateEducation in the service layer update a education of specific profile.
func (eduSvc *service) UpdateEducation(ctx context.Context, profileID int, eduID int, userID int, req specs.UpdateEducationRequest) (ID int, err error) {
	tx, _ := eduSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := eduSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	var value repository.UpdateEducationRepo
	value.Degree = req.Education.Degree
	value.UniversityName = req.Education.UniversityName
	value.Place = req.Education.Place
	value.PercentageOrCgpa = req.Education.PercentageOrCgpa
	value.PassingYear = req.Education.PassingYear
	value.UpdatedAt = today
	value.UpdatedByID = userID

	profileID, err = eduSvc.EducationRepo.UpdateEducation(ctx, profileID, eduID, value, tx)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("education(s) update with profile id : ", profileID)

	return profileID, nil
}
