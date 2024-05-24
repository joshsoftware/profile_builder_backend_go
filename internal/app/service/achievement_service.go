package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type achievementService struct {
	AchievementRepo repository.AchievementStorer
}

type AchievementService interface {
	CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error)
}

func NewAchivementService(db *pgx.Conn) AchievementService {
	achievementRepo := repository.NewAchievementRepo(db)

	return &achievementService{
		AchievementRepo: achievementRepo,
	}
}

// CreateAchievement : Service layer function adds certicates details to a user profile.
func (achSvc *achievementService) CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error) {
	if len(cDetail.Achievements) == 0 {
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
		return 0, err
	}

	return cDetail.ProfileID, nil
}