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

// ExperienceStore implements the ExperienceStorer interface.
type ExperienceStore struct {
	db *pgx.Conn
}

// ExperienceStorer defines methods to interact with user experience related data.
type ExperienceStorer interface {
	CreateExperience(ctx context.Context, values []ExperienceDao) error
}

// NewExperienceRepo creates a new instance of ExperienceRepo.
func NewExperienceRepo(db *pgx.Conn) ExperienceStorer {
	return &ExperienceStore{
		db: db,
	}
}

// CreateExperience inserts experience details into the database.
func (profileStore *ExperienceStore) CreateExperience(ctx context.Context, values []ExperienceDao) error {

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
	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
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
