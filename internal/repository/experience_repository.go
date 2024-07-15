package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// ExperienceStore implements the ExperienceStorer interface.
type ExperienceStore struct {
	db *pgxpool.Pool
}

// ExperienceStorer defines methods to interact with user experience related data.
type ExperienceStorer interface {
	CreateExperience(ctx context.Context, values []ExperienceRepo, tx pgx.Tx) error
	ListExperiences(ctx context.Context, profileID int, filter specs.ListExperiencesFilter, tx pgx.Tx) (values []specs.ExperienceResponse, err error)
	UpdateExperience(ctx context.Context, profileID int, eduID int, req UpdateExperienceRepo, tx pgx.Tx) (int, error)
	DeleteExperience(ctx context.Context, profileID, experienceID int, tx pgx.Tx) error
}

// NewExperienceRepo creates a new instance of ExperienceRepo.
func NewExperienceRepo(db *pgxpool.Pool) ExperienceStorer {
	return &ExperienceStore{
		db: db,
	}
}

// CreateExperience inserts experience details into the database.
func (expStore *ExperienceStore) CreateExperience(ctx context.Context, values []ExperienceRepo, tx pgx.Tx) error {

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
	_, err = tx.Exec(ctx, insertQuery, args...)
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

// ListExperiences returns a details experiences in the Database that are currently available for perticular ID
func (expStore *ExperienceStore) ListExperiences(ctx context.Context, profileID int, filter specs.ListExperiencesFilter, tx pgx.Tx) (values []specs.ExperienceResponse, err error) {
	queryBuilder := sq.Select(constants.ResponseExperiencesColumns...).From("experiences").Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar)

	if len(filter.ExperiencesIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.ExperiencesIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating get experiences query: ", err)
		return []specs.ExperienceResponse{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
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

	return values, nil
}

// UpdateExperience updates experience details into the database.
func (expStore *ExperienceStore) UpdateExperience(ctx context.Context, profileID int, eduID int, req UpdateExperienceRepo, tx pgx.Tx) (int, error) {
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

	res, err := tx.Exec(ctx, updateQuery, args...)
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

func (expStore *ExperienceStore) DeleteExperience(ctx context.Context, profileID, experienceID int, tx pgx.Tx) error {
	deleteQuery, args, err := psql.Delete("experiences").Where(sq.Eq{"id": experienceID, "profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating experience delete query: ", err)
		return err
	}

	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing experience delete query: ", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}

	return nil
}
