package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// ProfileStore implements the ProfileStorer interface.
type ProfileStore struct {
	db *pgxpool.Pool
}

// ProfileStorer defines methods to interact with user profile data.
type ProfileStorer interface {
	CreateProfile(ctx context.Context, pd ProfileRepo, tx pgx.Tx) (int, error)
	ListProfiles(ctx context.Context, tx pgx.Tx) (values []specs.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID int, tx pgx.Tx) (value specs.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo, tx pgx.Tx) (int, error)
	UpdateProfileStatus(ctx context.Context, profileID int, updateProfileStatus UpdateProfileStatusRepo, tx pgx.Tx) error
	DeleteProfile(ctx context.Context, profileID int, tx pgx.Tx) (err error)
	ListSkills(ctx context.Context, tx pgx.Tx) (values specs.ListSkills, err error)
	BeginTransaction(ctx context.Context) (tx pgx.Tx, err error)
	HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error)
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewProfileRepo(db *pgxpool.Pool) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// CreateProfile inserts a new user profile into the database.
func (profileStore *ProfileStore) CreateProfile(ctx context.Context, pd ProfileRepo, tx pgx.Tx) (int, error) {

	values := []interface{}{
		pd.Name, pd.Email, pd.Gender, pd.Mobile, pd.Designation, pd.Description, pd.Title,
		pd.YearsOfExperience, pd.PrimarySkills, pd.SecondarySkills, pd.GithubLink, pd.LinkedinLink, pd.CareerObjectives, 1, 1, pd.CreatedAt, pd.UpdatedAt, pd.CreatedByID, pd.UpdatedByID,
	}

	insertQuery, args, err := psql.Insert("profiles").
		Columns(constants.CreateUserColumns...).
		Values(values...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		zap.S().Error("Error generating profile insert query: ", err)
		return 0, err
	}

	if tx == nil {
		zap.S().Error("Transaction is nil")
		return 0, fmt.Errorf("internal error: transaction is nil")
	}

	var profileID int
	err = tx.QueryRow(ctx, insertQuery, args...).Scan(&profileID)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return 0, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing insert query: ", err)
		return 0, err
	}

	return profileID, nil
}

// ListProfiles returns a list of all profiles in the Database that are currently available
func (profileStore *ProfileStore) ListProfiles(ctx context.Context, tx pgx.Tx) (values []specs.ListProfiles, err error) {
	sql, args, err := psql.Select(constants.ListProfilesColumns...).From("profiles").OrderBy("created_at DESC").ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return []specs.ListProfiles{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing list project insert query:", err)
		return []specs.ListProfiles{}, err
	}

	for rows.Next() {
		var value specs.ListProfiles
		err := rows.Scan(&value.ID, &value.Name, &value.Email, &value.YearsOfExperience, &value.PrimarySkills, &value.IsCurrentEmployee, &value.IsActive)
		if err != nil {
			zap.S().Error("error scanning row:", err)
			return []specs.ListProfiles{}, err
		}
		values = append(values, value)
	}
	defer rows.Close()
	return values, nil
}

// ListSkills returns a list of all skills in the Database that are currently available
func (profileStore *ProfileStore) ListSkills(ctx context.Context, tx pgx.Tx) (values specs.ListSkills, err error) {
	sql, args, err := psql.Select("name").From("skills").ToSql()
	if err != nil {
		zap.S().Error("Error generating list skills select query: ", err)
		return specs.ListSkills{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing list skills insert query:", err)
		return specs.ListSkills{}, err
	}

	for rows.Next() {
		var value string
		rows.Scan(&value)
		values.Name = append(values.Name, value)
	}
	defer rows.Close()

	return values, nil
}

// GetProfile returns a details profile in the Database that are currently available for perticular ID
func (profileStore *ProfileStore) GetProfile(ctx context.Context, profileID int, tx pgx.Tx) (value specs.ResponseProfile, err error) {
	query, args, err := psql.Select(constants.ResponseProfileColumns...).From("profiles").
		Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return specs.ResponseProfile{}, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		zap.S().Error("Error executing get profile query: ", err)
		return specs.ResponseProfile{}, err
	}

	if rows.Next() {
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Email, &value.Gender, &value.Mobile, &value.Designation, &value.Description, &value.Title, &value.YearsOfExperience, &value.PrimarySkills, &value.SecondarySkills, &value.GithubLink, &value.LinkedinLink, &value.CareerObjectives); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return specs.ResponseProfile{}, err
		}
	} else {
		zap.S().Info("No profile found for profileID: ", profileID)
		return specs.ResponseProfile{}, errors.ErrNoRecordFound
	}
	defer rows.Close()

	return value, nil
}

// UpdateProfile updates an existing user profile in the database.
func (profileStore *ProfileStore) UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo, tx pgx.Tx) (int, error) {

	updateQuery, args, err := psql.Update("profiles").
		SetMap(map[string]interface{}{
			"name": pd.Name, "email": pd.Email,
			"gender": pd.Gender, "mobile": pd.Mobile,
			"designation": pd.Designation, "description": pd.Description,
			"title": pd.Title, "years_of_experience": pd.YearsOfExperience,
			"primary_skills": pd.PrimarySkills, "secondary_skills": pd.SecondarySkills,
			"github_link": pd.GithubLink, "linkedin_link": pd.LinkedinLink,
			"updated_at": pd.UpdatedAt, "updated_by_id": pd.UpdatedByID,
		}).
		Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating profile update query: ", err)
		return 0, err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return 0, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : profile")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}

func (profileStore *ProfileStore) DeleteProfile(ctx context.Context, profileID int, tx pgx.Tx) (err error) {
	deleteQuery, args, err := psql.Delete("profiles").Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID).Error("Error generating profile delete query: ", err)
		return err
	}

	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete profile query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}
	return nil
}

func (profileStore *ProfileStore) UpdateProfileStatus(ctx context.Context, profileID int, updateProfileStatus UpdateProfileStatusRepo, tx pgx.Tx) error {
	updateQuery := psql.Update("profiles")
	if updateProfileStatus.IsCurrentEmployee != "" {
		isCurrentEmployee := 0
		if updateProfileStatus.IsCurrentEmployee == "YES" {
			isCurrentEmployee = 1
		}
		updateQuery = updateQuery.Set("is_current_employee", isCurrentEmployee)
	}

	if updateProfileStatus.IsActive != "" {
		isActive := 0
		if updateProfileStatus.IsActive == "YES" {
			isActive = 1
		}
		updateQuery = updateQuery.Set("is_active", isActive)
	}

	updateQuery = updateQuery.Where(sq.Eq{"id": profileID})

	query, args, err := updateQuery.ToSql()

	fmt.Println("query in update profile status : ", query)
	fmt.Println("args in update profile status : ", args)
	fmt.Println("err in update profile status : ", err)
	if err != nil {
		zap.S().Error("Error generating update profile status query: ", err)
		return err
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return errors.ErrInvalidProfile
		}
		zap.S().Error("Error while executing update query: ", err)
		return err
	}
	fmt.Println("result in update profile status : ", res)

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : profile status")
		return errors.ErrInvalidRequestData
	}
	return nil
}

// BeginTransaction used to begin transaction while each task
func (profileStore *ProfileStore) BeginTransaction(ctx context.Context) (tx pgx.Tx, err error) {
	tx, err = profileStore.db.BeginTx(ctx, pgx.TxOptions{})
	return
}

// HandleTransaction used to handle transaction while each task
func (profileStore *ProfileStore) HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error) {
	if incomingErr != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return
		}
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		return
	}
	return
}
