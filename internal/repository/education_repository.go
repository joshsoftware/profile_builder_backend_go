package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
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
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewEducationRepo(db *pgx.Conn) EducationStorer {
	return &EducationStore{
		db: db,
	}
}

// CreateEducation inserts education details into the database.
func (profileStore *EducationStore) CreateEducation(ctx context.Context, values []EducationDao) error {

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
	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
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
