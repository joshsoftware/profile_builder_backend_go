package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// ExperienceStore implements the ExperienceStorer interface.
type ExperienceStore struct {
	db *pgx.Conn
}

// ExperienceStorer defines methods to interact with user experience related data.
type ExperienceStorer interface {
	CreateExperience(ctx context.Context, values []ExperienceDao) error
	GetExperiences(ctx context.Context, profileID int) (values []dto.ExperienceResponse, err error)
}

// NewExperienceRepo creates a new instance of ExperienceRepo.
func NewExperienceRepo(db *pgx.Conn) ExperienceStorer {
	return &ExperienceStore{
		db: db,
	}
}

// CreateExperience inserts experience details into the database.
func (expStore *ExperienceStore) CreateExperience(ctx context.Context, values []ExperienceDao) error {

	insertBuilder := sq.Insert("experiences").
		Columns(constants.CreateExperienceColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Designation, value.CompanyName, value.FromDate, value.ToDate,
			value.CreatedAt, value.UpdatedAt, value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating experience insert query: ", err)
		return err
	}
	_, err = expStore.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing create experience insert query: ", err)
		return err
	}

	return nil
}

// GetExperiences returns a details experiences in the Database that are currently available for perticular ID
func (expStore *ExperienceStore) GetExperiences(ctx context.Context, profileID int) (values []dto.ExperienceResponse, err error) {
	sql, args, err := sq.Select(constants.ResponseExperiencesColumns...).From("experiences").
		Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating get experiences select query: ", err)
		return []dto.ExperienceResponse{}, err
	}

	rows, err := expStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get experiences query: ", err)
		return []dto.ExperienceResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value dto.ExperienceResponse
		if err := rows.Scan(&value.ProfileID, &value.Designation, &value.CompanyName, &value.FromDate, &value.ToDate); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []dto.ExperienceResponse{}, err
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		zap.S().Info("No experience found for profileID: ", profileID)
		return []dto.ExperienceResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}
