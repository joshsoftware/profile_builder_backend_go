package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// AchievementService represents a set of methods for accessing the achievements
type AchievementService interface {
	CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, ID int) (profileID int, err error)
	UpdateAchievement(ctx context.Context, profileID string, eduID string, req specs.UpdateAchievementRequest) (id int, err error)
}

// CreateAchievement : Service layer function adds certicates details to a user profile.
func (achSvc *service) CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, id int) (profileID int, err error) {
	if len(cDetail.Achievements) == 0 {
		zap.S().Error("achievements payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.AchievementRepo
	for i := 0; i < len(cDetail.Achievements); i++ {
		var val repository.AchievementRepo

		val.ProfileID = id
		val.Name = cDetail.Achievements[i].Name
		val.Description = cDetail.Achievements[i].Description
		val.CreatedAt = today
		val.UpdatedAt = today
		//TODO by context
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = achSvc.AchievementRepo.CreateAchievement(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create achievement : ", err, "for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("achievement(s) created with profile id : ", id)

	return id, nil
}

// UpdateAchievement in the service layer update a achievements of specific profile.
func (achSvc *service) UpdateAchievement(ctx context.Context, profileID string, eduID string, req specs.UpdateAchievementRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get achievement params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	var value repository.UpdateAchievementRepo
	value.Name = req.Achievement.Name
	value.Description = req.Achievement.Description
	value.UpdatedAt = today
	//TODO by context
	value.UpdatedByID = 1

	id, err = achSvc.AchievementRepo.UpdateAchievement(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update achievement : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("achievement(s) update with profile id : ", id)

	return id, nil
}
