package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
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
	ListAchievements(ctx context.Context, profileID int, filter dto.ListAchievementFilter) ([]dto.AchievementResponse, error)
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

func (achStore *AchievementStore) ListAchievements(ctx context.Context, profileID int, filter dto.ListAchievementFilter) (values []dto.AchievementResponse, err error) {
	queryBuilder := sq.Select(constants.ResponseAchievementsColumns...).From("achievements").Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar)

	if len(filter.AchievementIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.AchievementIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating get achievements query: ", err)
		return []dto.AchievementResponse{}, err
	}
	fmt.Println("sql , args  : ", sql, args)
	rows, err := achStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get achievements query: ", err)
		return []dto.AchievementResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var val dto.AchievementResponse
		err = rows.Scan(&val.ID, &val.ProfileID, &val.Name, &val.Description)
		if err != nil {
			zap.S().Error("Error scanning achievements rows: ", err)
			return []dto.AchievementResponse{}, err
		}
		values = append(values, val)
	}

	if len(values) == 0 {
		zap.S().Info("No achievements found for profileID: ", profileID)
		return []dto.AchievementResponse{}, nil
	}

	return values, nil
}
