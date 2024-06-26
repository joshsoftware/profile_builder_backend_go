package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// ExperienceStore implements the ExperienceStorer interface.
type ExperienceStore struct {
	db *pgx.Conn
}

// ExperienceStorer defines methods to interact with user experience related data.
type ExperienceStorer interface {
	CreateExperience(ctx context.Context, values []ExperienceRepo) error
	GetExperiences(ctx context.Context, profileID int) (values []specs.ExperienceResponse, err error)
	UpdateExperience(ctx context.Context, profileID int, eduID int, req UpdateExperienceRepo) (int, error)
}

// NewExperienceRepo creates a new instance of ExperienceRepo.
func NewExperienceRepo(db *pgx.Conn) ExperienceStorer {
	return &ExperienceStore{
		db: db,
	}
}

// CreateExperience inserts experience details into the database.
func (expStore *ExperienceStore) CreateExperience(ctx context.Context, values []ExperienceRepo) error {

	insertBuilder := psql.Insert("experiences").
		Columns(constants.CreateExperienceColumns...)

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
func (expStore *ExperienceStore) GetExperiences(ctx context.Context, profileID int) (values []specs.ExperienceResponse, err error) {
	sql, args, err := psql.Select(constants.ResponseExperiencesColumns...).From("experiences").
		Where(sq.Eq{"profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating get experiences select query: ", err)
		return []specs.ExperienceResponse{}, err
	}

	rows, err := expStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get experiences query: ", err)
		return []specs.ExperienceResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value specs.ExperienceResponse
		if err := rows.Scan(&value.ID, &value.ProfileID, &value.Designation, &value.CompanyName, &value.FromDate, &value.ToDate); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []specs.ExperienceResponse{}, err
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		zap.S().Info("No experience found for profileID: ", profileID)
		return []specs.ExperienceResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}

// UpdateExperience updates experience details into the database.
func (expStore *ExperienceStore) UpdateExperience(ctx context.Context, profileID int, eduID int, req UpdateExperienceRepo) (int, error) {
	updateQuery, args, err := psql.Update("experiences").
		SetMap(map[string]interface{}{
			"designation": req.Designation, "company_name": req.CompanyName,
			"from_date": req.FromDate, "to_date": req.ToDate,
			"updated_at": req.UpdatedAt, "updated_by_id": req.UpdatedByID,
		}).Where(sq.Eq{"profile_id": profileID, "id": eduID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating experience update query: ", err)
		return 0, err
	}

	res, err := expStore.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing experience update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : experience")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}
