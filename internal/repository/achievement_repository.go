package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// AchievementStore implements the AchievementStorer interface.
type AchievementStore struct {
	db *pgx.Conn
}

// NewAchievementRepo creates a new instance of AchievementRepo.
func NewAchievementRepo(db *pgx.Conn) AchievementStorer {
	return &AchievementStore{
		db: db,
	}
}

// AchievementStorer defines methods to interact with user achievement related data.
type AchievementStorer interface {
	CreateAchievement(ctx context.Context, values []AchievementDao) error
}

// CreateAchievement inserts achievements details into the database.
func (profileStore *AchievementStore) CreateAchievement(ctx context.Context, values []AchievementDao) error {

	insertBuilder := sq.Insert("achievements").
		Columns(constants.CreateAchievementColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.Description, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating achievement insert query: ", err)
		return err
	}
	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing create achievement insert query: ", err)
		return err
	}

	return nil
}
