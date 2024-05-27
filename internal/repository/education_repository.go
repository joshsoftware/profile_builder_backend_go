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

// EducationStore implements the EducationStorer interface.
type EducationStore struct {
	db *pgx.Conn
}

// EducationStorer defines methods to interact with user education related data.
type EducationStorer interface {
	CreateEducation(ctx context.Context, values []EducationDao) error
	GetEducation(ctx context.Context, profileID int) (value []dto.EducationResponse, err error)
	UpdateEducation(ctx context.Context, profileID int, eduID int, req dto.UpdateEducationRequest) (int, error)
}

// NewEducationRepo creates a new instance of EducationRepo.
func NewEducationRepo(db *pgx.Conn) EducationStorer {
	return &EducationStore{
		db: db,
	}
}

// CreateEducation inserts education details into the database.
func (eduStore *EducationStore) CreateEducation(ctx context.Context, values []EducationDao) error {

	insertBuilder := sq.Insert("educations").
		Columns(constants.CreateEducationColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Degree, value.UniversityName, value.Place, value.PercentageOrCgpa,
			value.PassingYear, value.CreatedAt, value.UpdatedAt, value.CreatedByID,
			value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating education insert query: ", err)
		return err
	}
	_, err = eduStore.db.Exec(ctx, insertQuery, args...)
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

// GetEducation returns a details education in the Database that are currently available for perticular ID
func (eduStore *EducationStore) GetEducation(ctx context.Context, profileID int) (values []dto.EducationResponse, err error) {
	sql, args, err := sq.Select(constants.ResponseEducationColumns...).From("educations").
		Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating get education select query: ", err)
		return []dto.EducationResponse{}, err
	}

	rows, err := eduStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get education query: ", err)
		return []dto.EducationResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value dto.EducationResponse
		if err := rows.Scan(&value.ProfileID, &value.Degree, &value.UniversityName, &value.Place,
			&value.PercentageOrCgpa, &value.PassingYear); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []dto.EducationResponse{}, err
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		zap.S().Info("No education found for profileID: ", profileID)
		return []dto.EducationResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}

// UpdateEducation updates education details into the database.
func (eduStore *EducationStore) UpdateEducation(ctx context.Context, profileID int, eduID int, req dto.UpdateEducationRequest) (int, error) {
	updateQuery, args, err := sq.Update("educations").
		Set("degree", req.Education.Degree).
		Set("university_name", req.Education.UniversityName).
		Set("place", req.Education.Place).
		Set("percent_or_cgpa", req.Education.PercentageOrCgpa).
		Set("passing_year", req.Education.PassingYear).
		Where(sq.Eq{"profile_id": profileID, "id": eduID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating education update query: ", err)
		return 0, err
	}

	res, err := eduStore.db.Exec(ctx, updateQuery, args...)
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
