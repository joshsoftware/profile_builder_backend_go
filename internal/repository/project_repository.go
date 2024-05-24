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

// ProfileStore implements the ProfileStorer interface.
type ProjectStore struct {
	db *pgx.Conn
}

// ProfileStorer defines methods to interact with user profile data.
type ProjectStorer interface {
	CreateProject(ctx context.Context, values []ProjectDao) error
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewProjectRepo(db *pgx.Conn) ProjectStorer {
	return &ProjectStore{
		db: db,
	}
}

// CreateProject inserts project details into the database.
func (profileStore *ProjectStore) CreateProject(ctx context.Context, values []ProjectDao) error {

	insertBuilder := sq.Insert("projects").
		Columns(constants.CreateProjectColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.Description, value.Role, value.Responsibilities,
			value.Technologies, value.TechWorkedOn, value.WorkingStartDate,
			value.WorkingEndDate, value.Duration, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating project insert query: ", err)
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
		zap.S().Error("Error executing create project insert query: ", err)
		return err
	}

	return nil
}
