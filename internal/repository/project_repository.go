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

// ProjectStore implements the ProfileStorer interface.
type ProjectStore struct {
	db *pgx.Conn
}

// ProjectStorer defines methods to interact with user profile data.
type ProjectStorer interface {
	CreateProject(ctx context.Context, values []ProjectDao) error
	GetProjects(ctx context.Context, profileID int) (values []dto.ProjectResponse, err error)
}

// NewProjectRepo creates a new instance of ProfileRepo.
func NewProjectRepo(db *pgx.Conn) ProjectStorer {
	return &ProjectStore{
		db: db,
	}
}

// CreateProject inserts project details into the database.
func (projectStore *ProjectStore) CreateProject(ctx context.Context, values []ProjectDao) error {

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
	_, err = projectStore.db.Exec(ctx, insertQuery, args...)
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

// GetProjects returns a details projects in the Database that are currently available for perticular ID
func (projectStore *ProjectStore) GetProjects(ctx context.Context, profileID int) (values []dto.ProjectResponse, err error) {
	sql, args, err := sq.Select(constants.ResponseProjectsColumns...).From("projects").
		Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating get projects select query: ", err)
		return []dto.ProjectResponse{}, err
	}

	rows, err := projectStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get projects query: ", err)
		return []dto.ProjectResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value dto.ProjectResponse
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Description, &value.Role, &value.Responsibilities, &value.Technologies, &value.TechWorkedOn, &value.WorkingStartDate,
			&value.WorkingEndDate, &value.Duration); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []dto.ProjectResponse{}, err
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		zap.S().Info("No project found for profileID: ", profileID)
		return []dto.ProjectResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}
