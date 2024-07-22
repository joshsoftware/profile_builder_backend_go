package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// AchievementStore implements the AchievementStorer interface.
type AchievementStore struct {
	db *pgxpool.Pool
}

// NewAchievementRepo creates a new instance of AchievementRepo.
func NewAchievementRepo(db *pgxpool.Pool) AchievementStorer {
	return &AchievementStore{
		db: db,
	}
}

// AchievementStorer defines methods to interact with user achievement related data.
type AchievementStorer interface {
	CreateAchievement(ctx context.Context, values []AchievementRepo, tx pgx.Tx) error
	UpdateAchievement(ctx context.Context, profileID int, achID int, req UpdateAchievementRepo, tx pgx.Tx) (int, error)
	ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter, tx pgx.Tx) ([]specs.AchievementResponse, error)
	DeleteAchievement(ctx context.Context, profileID, achievementID int, tx pgx.Tx) error
}

// CreateAchievement inserts achievements details into the database.
func (achStore *AchievementStore) CreateAchievement(ctx context.Context, values []AchievementRepo, tx pgx.Tx) error {

	insertBuilder := psql.Insert("achievements").
		Columns(constants.CreateAchievementColumns...)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.Description, value.Priorities, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}
	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating achievement insert query: ", err)
		return err
	}

	_, err = tx.Exec(ctx, insertQuery, args...)
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
func (achStore *AchievementStore) UpdateAchievement(ctx context.Context, profileID int, achID int, req UpdateAchievementRepo, tx pgx.Tx) (int, error) {
	updateQuery, args, err := psql.Update("achievements").
		SetMap(map[string]interface{}{
			"name": req.Name, "description": req.Description,
			"updated_at": req.UpdatedAt, "updated_by_id": req.UpdatedByID,
		}).Where(sq.Eq{"profile_id": profileID, "id": achID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating achievement update query: ", err)
		return 0, err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
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

// ListAchievements lists achievements from the database.
func (achStore *AchievementStore) ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter, tx pgx.Tx) (values []specs.AchievementResponse, err error) {
	queryBuilder := psql.Select(constants.ResponseAchievementsColumns...).From("achievements").Where(sq.Eq{"profile_id": profileID}).OrderBy("priorities")

	if len(filter.AchievementIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.AchievementIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating get achievements query: ", err)
		return []specs.AchievementResponse{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get achievements query: ", err)
		return []specs.AchievementResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var val specs.AchievementResponse
		err = rows.Scan(&val.ID, &val.ProfileID, &val.Name, &val.Description)
		if err != nil {
			zap.S().Error("Error scanning achievements rows: ", err)
			return []specs.AchievementResponse{}, err
		}
		values = append(values, val)
	}

	return values, nil
}

// DeleteAchievement deletes achievements details into the database.
func (achStore *AchievementStore) DeleteAchievement(ctx context.Context, profileID, achievementID int, tx pgx.Tx) error {
	deleteQuery, args, err := psql.Delete("achievements").Where(sq.Eq{"id": achievementID, "profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID, "achievement_id ", achievementID).Error("Error generating delete achievement query: ", zap.Error(err))
		return err
	}
	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete achievement query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}
	return nil
}
