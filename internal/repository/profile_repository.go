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

// ProfileStore implements the ProfileStorer interface.
type ProfileStore struct {
	db *pgx.Conn
}

// ProfileStorer defines methods to interact with user profile data.
type ProfileStorer interface {
	CreateProfile(ctx context.Context, pd ProfileRepo) (int, error)
	ListProfiles(ctx context.Context) (values []specs.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID int) (value specs.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo) (int, error)
	ListSkills(ctx context.Context) (values specs.ListSkills, err error)
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewProfileRepo(db *pgx.Conn) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// CreateProfile inserts a new user profile into the database.
func (profileStore *ProfileStore) CreateProfile(ctx context.Context, pd ProfileRepo) (int, error) {

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

	var profileID int
	err = profileStore.db.QueryRow(ctx, insertQuery, args...).Scan(&profileID)
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
func (profileStore *ProfileStore) ListProfiles(ctx context.Context) (values []specs.ListProfiles, err error) {
	sql, args, err := psql.Select(constants.ListProfilesColumns...).From("profiles").ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return []specs.ListProfiles{}, err
	}
	rows, err := profileStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing list project insert query:", err)
		return []specs.ListProfiles{}, err
	}

	for rows.Next() {
		var value specs.ListProfiles
		err := rows.Scan(&value.ID, &value.Name, &value.Email, &value.YearsOfExperience, &value.PrimarySkills, &value.IsCurrentEmployee)
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
func (profileStore *ProfileStore) ListSkills(ctx context.Context) (values specs.ListSkills, err error) {
	sql, args, err := psql.Select("name").From("skills").ToSql()
	if err != nil {
		zap.S().Error("Error generating list skills select query: ", err)
		return specs.ListSkills{}, err
	}
	rows, err := profileStore.db.Query(ctx, sql, args...)
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
func (profileStore *ProfileStore) GetProfile(ctx context.Context, profileID int) (value specs.ResponseProfile, err error) {
	query, args, err := psql.Select(constants.ResponseProfileColumns...).From("profiles").
		Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return specs.ResponseProfile{}, err
	}

	rows, err := profileStore.db.Query(ctx, query, args...)
	if err != nil {
		zap.S().Error("Error executing get profile query: ", err)
		return specs.ResponseProfile{}, err
	}

	if rows.Next() {
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Email, &value.Gender, &value.Mobile, &value.Designation, &value.Description, &value.Title, &value.YearsOfExperience, &value.PrimarySkills, &value.SecondarySkills, &value.GithubLink, &value.LinkedinLink); err != nil {
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
func (profileStore *ProfileStore) UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo) (int, error) {

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

	res, err := profileStore.db.Exec(ctx, updateQuery, args...)
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
