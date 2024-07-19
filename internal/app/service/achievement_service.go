package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// AchievementService represents a set of methods for accessing the achievements
type AchievementService interface {
	CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, profileID int, userID int) (ID int, err error)
	UpdateAchievement(ctx context.Context, profileID int, achID int, userID int, req specs.UpdateAchievementRequest) (ID int, err error)
	ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter) (value []specs.AchievementResponse, err error)
	DeleteAchievement(ctx context.Context, profileID, achievementID int) error
}

// CreateAchievement : Service layer function adds achievement details to a user profile.
func (achSvc *service) CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, profileID int, userID int) (ID int, err error) {
	tx, _ := achSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := achSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	count, err := achSvc.ProfileRepo.CountRecords(ctx, profileID, constants.Achievements, tx)
	if err != nil {
		return 0, errors.ErrInvalidRequestData
	}

	count++
	var value []repository.AchievementRepo
	for _, achievement := range cDetail.Achievements {
		val := repository.AchievementRepo{
			ProfileID:   profileID,
			Name:        achievement.Name,
			Description: achievement.Description,
			Priorities:  count,
			CreatedAt:   today,
			UpdatedAt:   today,
			CreatedByID: userID,
			UpdatedByID: userID,
		}

		count++
		value = append(value, val)
	}

	err = achSvc.AchievementRepo.CreateAchievement(ctx, value, tx)
	if err != nil {
		zap.S().Error("Unable to create achievement : ", err, "for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("achievement(s) added with profile id : ", profileID)

	return profileID, nil
}

// UpdateAchievement in the service layer update a achievements of specific profile.
func (achSvc *service) UpdateAchievement(ctx context.Context, profileID int, achID int, userID int, req specs.UpdateAchievementRequest) (ID int, err error) {
	tx, _ := achSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := achSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	var value repository.UpdateAchievementRepo
	value.Name = req.Achievement.Name
	value.Description = req.Achievement.Description
	value.UpdatedAt = today
	value.UpdatedByID = userID

	profileID, err = achSvc.AchievementRepo.UpdateAchievement(ctx, profileID, achID, value, tx)
	if err != nil {
		zap.S().Error("Unable to update achievement : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("achievement(s) update with profile id : ", profileID)

	return profileID, nil
}

func (achSvc *service) ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter) (value []specs.AchievementResponse, err error) {
	tx, _ := achSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := achSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	value, err = achSvc.AchievementRepo.ListAchievements(ctx, profileID, filter, tx)
	if err != nil {
		zap.S().Error("error to get achievement : ", err, " for profile id : ", profileID)
		return []specs.AchievementResponse{}, err
	}
	return value, nil
}

func (achSvc *service) DeleteAchievement(ctx context.Context, profileID, achievementID int) (err error) {
	tx, _ := achSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := achSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	err = achSvc.AchievementRepo.DeleteAchievement(ctx, profileID, achievementID, tx)

	if err != nil {
		if err == errors.ErrNoData {
			zap.S().Warn("No achievement found to delete for achievement id: ", achievementID, " for profile id: ", profileID)
			return err
		}
		zap.S().Error("Error deleting achievement: ", err, " for achievement id: ", achievementID, " for profile id: ", profileID)
		return err
	}
	zap.S().Info("Achievement deleted successfully for achievement id: ", achievementID, " for profile id: ", profileID)
	return nil
}
