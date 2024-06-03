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
	CreateAchievement(ctx context.Context, values []AchievementRepo) error
	UpdateAchievement(ctx context.Context, profileID int, achID int, req UpdateAchievementRepo) (int, error)
}

// CreateAchievement inserts achievements details into the database.
func (achStore *AchievementStore) CreateAchievement(ctx context.Context, values []AchievementRepo) error {

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
	_, err = achStore.db.Exec(ctx, insertQuery, args...)
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

// UpdateAchievement updates achievements details into the database.
func (achStore *AchievementStore) UpdateAchievement(ctx context.Context, profileID int, achID int, req UpdateAchievementRepo) (int, error) {
	updateQuery, args, err := sq.Update("achievements").Set("name", req.Name).Set("description", req.Description).Set("updated_at", req.UpdatedAt).Set("updated_by_id", req.UpdatedByID).Where(sq.Eq{"profile_id": profileID, "id": achID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating achievement update query: ", err)
		return 0, err
	}

	res, err := achStore.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing achievement update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : achievement")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}
