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

// ProjectStore implements the ProfileStorer interface.
type ProjectStore struct {
	db *pgx.Conn
}

// ProjectStorer defines methods to interact with user profile data.
type ProjectStorer interface {
	CreateProject(ctx context.Context, values []ProjectRepo) error
	GetProjects(ctx context.Context, profileID int) (values []specs.ProjectResponse, err error)
	UpdateProject(ctx context.Context, profileID int, eduID int, req UpdateProjectRepo) (int, error)
}

// NewProjectRepo creates a new instance of ProfileRepo.
func NewProjectRepo(db *pgx.Conn) ProjectStorer {
	return &ProjectStore{
		db: db,
	}
}

// CreateProject inserts project details into the database.
func (projectStore *ProjectStore) CreateProject(ctx context.Context, values []ProjectRepo) error {

	insertBuilder := psql.Insert("projects").
		Columns(constants.CreateProjectColumns...)

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
func (projectStore *ProjectStore) GetProjects(ctx context.Context, profileID int) (values []specs.ProjectResponse, err error) {
	sql, args, err := psql.Select(constants.ResponseProjectsColumns...).From("projects").
		Where(sq.Eq{"profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating get projects select query: ", err)
		return []specs.ProjectResponse{}, err
	}

	rows, err := projectStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get projects query: ", err)
		return []specs.ProjectResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value specs.ProjectResponse
		if err := rows.Scan(&value.ID, &value.ProfileID, &value.Name, &value.Description, &value.Role, &value.Responsibilities, &value.Technologies, &value.TechWorkedOn, &value.WorkingStartDate,
			&value.WorkingEndDate, &value.Duration); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return []specs.ProjectResponse{}, err
		}
		values = append(values, value)
	}

	if len(values) == 0 {
		zap.S().Info("No project found for profileID: ", profileID)
		return []specs.ProjectResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}

// UpdateProject updates projects details into the database.
func (projectStore *ProjectStore) UpdateProject(ctx context.Context, profileID int, eduID int, req UpdateProjectRepo) (int, error) {
	updateQuery, args, err := psql.Update("projects").
		SetMap(map[string]interface{}{
			"name": req.Name, "description": req.Description,
			"role": req.Role, "responsibilities": req.Responsibilities,
			"technologies": req.Technologies, "tech_worked_on": req.TechWorkedOn,
			"working_start_date": req.WorkingStartDate, "working_end_date": req.WorkingEndDate,
			"duration": req.Duration, "updated_at": req.UpdatedAt,
			"updated_by_id": req.UpdatedByID,
		}).
		Where(sq.Eq{"profile_id": profileID, "id": eduID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating projects update query: ", err)
		return 0, err
	}

	res, err := projectStore.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing projects update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : projects")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}
