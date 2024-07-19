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

// ProjectStore implements the ProfileStorer interface.
type ProjectStore struct {
	db *pgxpool.Pool
}

// ProjectStorer defines methods to interact with user profile data.
type ProjectStorer interface {
	CreateProject(ctx context.Context, values []ProjectRepo, tx pgx.Tx) error
	ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter, tx pgx.Tx) (values []specs.ProjectResponse, err error)
	UpdateProject(ctx context.Context, profileID int, eduID int, req UpdateProjectRepo, tx pgx.Tx) (int, error)
	DeleteProject(ctx context.Context, profileID, projectID int, tx pgx.Tx) error
}

// NewProjectRepo creates a new instance of ProfileRepo.
func NewProjectRepo(db *pgxpool.Pool) ProjectStorer {
	return &ProjectStore{
		db: db,
	}
}

// CreateProject inserts project details into the database.
func (projectStore *ProjectStore) CreateProject(ctx context.Context, values []ProjectRepo, tx pgx.Tx) error {

	insertBuilder := psql.Insert("projects").
		Columns(constants.CreateProjectColumns...)
	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.Description, value.Role, value.Responsibilities,
			value.Technologies, value.TechWorkedOn, value.WorkingStartDate,
			value.WorkingEndDate, value.Duration, value.Priorities, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating project insert query: ", err)
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
		zap.S().Error("Error executing create project insert query: ", err)
		return err
	}

	return nil
}

// ListProjects returns a details projects in the Database that are currently available for perticular profile ID
func (projectStore *ProjectStore) ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter, tx pgx.Tx) (values []specs.ProjectResponse, err error) {

	queryBuilder := psql.Select(constants.ResponseProjectsColumns...).From("projects").Where(sq.Eq{"profile_id": profileID}).OrderBy("priorities")

	if len(filter.ProjectsIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.ProjectsIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating get projects query: ", err)
		return []specs.ProjectResponse{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
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

	return values, nil
}

// UpdateProject updates projects details into the database.
func (projectStore *ProjectStore) UpdateProject(ctx context.Context, profileID int, eduID int, req UpdateProjectRepo, tx pgx.Tx) (int, error) {
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

	res, err := tx.Exec(ctx, updateQuery, args...)
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

// DeleteProject deletes projects details into the database.
func (projectStore *ProjectStore) DeleteProject(ctx context.Context, profileID, projectID int, tx pgx.Tx) error {
	deleteQuery, args, err := psql.Delete("projects").Where(sq.Eq{"id": projectID, "profile_id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID, "project_id ", projectID).Error("Error generating delete project query: ", zap.Error(err))
		return err
	}

	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete project query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}
	return nil
}
