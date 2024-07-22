package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// ExperienceService represents a set of methods for accessing the experiences
type ExperienceService interface {
	CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, profileID int, userID int) (ID int, err error)
	ListExperiences(ctx context.Context, id int, filter specs.ListExperiencesFilter) (values []specs.ExperienceResponse, err error)
	UpdateExperience(ctx context.Context, profileID int, expID int, userID int, req specs.UpdateExperienceRequest) (ID int, err error)
	DeleteExperience(ctx context.Context, profileID, experienceID int) error
}

// CreateExperience : Service layer function adds experiences details to a user profile.
func (expSvc *service) CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, profileID int, userID int) (ID int, err error) {
	tx, _ := expSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := expSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	count, err := expSvc.ProfileRepo.CountRecords(ctx, profileID, constants.Experiences, tx)
	if err != nil {
		return 0, errors.ErrInvalidRequestData
	}

	count++
	var value []repository.ExperienceRepo
	for _, experience := range expDetail.Experiences {
		val := repository.ExperienceRepo{
			ProfileID:   profileID,
			Designation: experience.Designation,
			CompanyName: experience.CompanyName,
			FromDate:    experience.FromDate,
			ToDate:      experience.ToDate,
			Priorities:  count,
			CreatedAt:   today,
			UpdatedAt:   today,
			CreatedByID: userID,
			UpdatedByID: userID,
		}

		count++
		value = append(value, val)
	}

	err = expSvc.ExperienceRepo.CreateExperience(ctx, value, tx)
	if err != nil {
		zap.S().Error("Unable to create experiences : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("experience(s) created with profile id : ", profileID)

	return profileID, nil
}

// ListExperiences in the service layer retrieves a experiences of specific profile.
func (expSvc *service) ListExperiences(ctx context.Context, id int, filter specs.ListExperiencesFilter) (values []specs.ExperienceResponse, err error) {
	tx, _ := expSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := expSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	values, err = expSvc.ExperienceRepo.ListExperiences(ctx, id, filter, tx)
	if err != nil {
		zap.S().Error("Unable to get experiences : ", err, " for profile id : ", id)
		return []specs.ExperienceResponse{}, err
	}
	return values, nil
}

// UpdateExperience in the service layer update a experience of specific profile.
func (expSvc *service) UpdateExperience(ctx context.Context, profileID int, expID int, userID int, req specs.UpdateExperienceRequest) (ID int, err error) {
	tx, _ := expSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := expSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	var value repository.UpdateExperienceRepo
	value.Designation = req.Experience.Designation
	value.CompanyName = req.Experience.CompanyName
	value.FromDate = req.Experience.FromDate
	value.ToDate = req.Experience.ToDate
	value.UpdatedAt = today
	value.UpdatedByID = userID

	profileID, err = expSvc.ExperienceRepo.UpdateExperience(ctx, profileID, expID, value, tx)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("experience(s) update with profile id : ", profileID)

	return profileID, nil
}

func (expSvc *service) DeleteExperience(ctx context.Context, profileID, experienceID int) (err error) {
	tx, _ := expSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := expSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	err = expSvc.ExperienceRepo.DeleteExperience(ctx, profileID, experienceID, tx)
	if err != nil {
		if err == errors.ErrNoData {
			zap.S().Warn("No experience found to delete for experience id: ", experienceID, " for profile id: ", profileID)
			return err
		}
		zap.S().Error("Error deleting experience: ", err, " for experience id: ", experienceID, " for profile id: ", profileID)
		return err
	}
	zap.S().Info("experience deleted successfully for experience id: ", experienceID, " for profile id: ", profileID)
	return nil
}
