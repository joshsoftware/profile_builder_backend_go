package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// AchievementService represents a set of methods for accessing the achievements
type AchievementService interface {
	CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error)
	ListAchievements(ctx context.Context, profileID int) (value []dto.AchievementResponse, err error)
}

// CreateAchievement : Service layer function adds certicates details to a user profile.
func (achSvc *service) CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error) {
	if len(cDetail.Achievements) == 0 {
		zap.S().Error("achievements payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}
	var value []repository.AchievementDao
	for i := 0; i < len(cDetail.Achievements); i++ {
		var val repository.AchievementDao

		val.ProfileID = cDetail.ProfileID
		val.Name = cDetail.Achievements[i].Name
		val.Description = cDetail.Achievements[i].Description
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = achSvc.AchievementRepo.CreateAchievement(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create achievement : ", err, "for profile id : ", profileID)
		return 0, err
	}

	return cDetail.ProfileID, nil
}

func (achSvc *service) ListAchievements(ctx context.Context, profileID int) (value []dto.AchievementResponse, err error) {
	value, err = achSvc.AchievementRepo.ListAchievements(ctx, profileID)
	if err != nil {
		zap.S().Error("error to get achievement : ", err, " for profile id : ", profileID)
		return []dto.AchievementResponse{}, err
	}

	return value, nil
}
