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

// EducationStore implements the EducationStorer interface.
type EducationStore struct {
	db *pgxpool.Pool
}

// EducationStorer defines methods to interact with user education related data.
type EducationStorer interface {
	CreateEducation(ctx context.Context, values []EducationRepo, tx pgx.Tx) error
	ListEducations(ctx context.Context, profileID int, filter specs.ListEducationsFilter, tx pgx.Tx) (values []specs.EducationResponse, err error)
	UpdateEducation(ctx context.Context, profileID int, eduID int, req UpdateEducationRepo, tx pgx.Tx) (int, error)
	DeleteEducation(ctx context.Context, profileID, educationID int, tx pgx.Tx) error
}

// NewEducationRepo creates a new instance of EducationRepo.
func NewEducationRepo(db *pgxpool.Pool) EducationStorer {
	return &EducationStore{
		db: db,
	}
}

// CreateEducation inserts education details into the database.
func (eduStore *EducationStore) CreateEducation(ctx context.Context, values []EducationRepo, tx pgx.Tx) error {

	insertBuilder := psql.Insert("educations").
		Columns(constants.CreateEducationColumns...)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Degree, value.UniversityName, value.Place, value.PercentageOrCgpa,
			value.PassingYear, value.Priorities, value.CreatedAt, value.UpdatedAt, value.CreatedByID,
			value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating education insert query: ", err)
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
		zap.S().Error("error executing create education insert query:", err)
		return err
	}

	return nil
}

// ListEducations returns a details education in the Database that are currently available for perticular ID
func (eduStore *EducationStore) ListEducations(ctx context.Context, profileID int, filter specs.ListEducationsFilter, tx pgx.Tx) (values []specs.EducationResponse, err error) {
	queryBuilder := psql.Select(constants.ResponseEducationColumns...).From("educations").Where(sq.Eq{"profile_id": profileID}).OrderBy("priorities")

	if len(filter.EduationsIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.EduationsIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating get educations query: ", err)
		return []specs.EducationResponse{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get educations query: ", err)
		return []specs.EducationResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value specs.EducationResponse
		if err := rows.Scan(&value.ProfileID, &value.ID, &value.Degree, &value.UniversityName, &value.Place, &value.PercentageOrCgpa, &value.PassingYear); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []specs.EducationResponse{}, err
		}
		values = append(values, value)
	}

	return values, nil
}

// UpdateEducation updates education details into the database.
func (eduStore *EducationStore) UpdateEducation(ctx context.Context, profileID int, eduID int, req UpdateEducationRepo, tx pgx.Tx) (int, error) {
	updateQuery, args, err := psql.Update("educations").
		SetMap(map[string]interface{}{
			"degree": req.Degree, "university_name": req.UniversityName,
			"place": req.Place, "percent_or_cgpa": req.PercentageOrCgpa,
			"passing_year": req.PassingYear, "updated_at": req.UpdatedAt,
			"updated_by_id": req.UpdatedByID,
		}).Where(sq.Eq{"profile_id": profileID, "id": eduID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating education update query: ", err)
		return 0, err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing education update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : education")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}

// DeleteEducation deletes education details into the database.
func (eduStore *EducationStore) DeleteEducation(ctx context.Context, profileID, educationID int, tx pgx.Tx) error {
	deleteQuery, args, err := psql.Delete("educations").Where(sq.Eq{"id": educationID, "profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID, "education_id ; ", educationID).Error("Error generating delete education query : ", zap.Error(err))
		return err
	}

	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete education query : ", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}
	return nil
}
