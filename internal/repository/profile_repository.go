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

// ProfileStore implements the ProfileStorer interface.
type ProfileStore struct {
	db *pgx.Conn
}

// ProfileStorer defines methods to interact with user profile data.
type ProfileStorer interface {
	CreateProfile(ctx context.Context, pd dto.CreateProfileRequest) (int, error)
	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID int) (value dto.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID int, pd dto.UpdateProfileRequest) (int, error)
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewProfileRepo(db *pgx.Conn) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

// CreateProfile inserts a new user profile into the database.
func (profileStore *ProfileStore) CreateProfile(ctx context.Context, pd dto.CreateProfileRequest) (int, error) {

	values := []interface{}{
		pd.Profile.Name, pd.Profile.Email, pd.Profile.Gender, pd.Profile.Mobile,
		pd.Profile.Designation, pd.Profile.Description, pd.Profile.Title,
		pd.Profile.YearsOfExperience, pd.Profile.PrimarySkills, pd.Profile.SecondarySkills,
		pd.Profile.GithubLink, pd.Profile.LinkedinLink, 1, 1,
	}

	insertQuery, args, err := sq.Insert("profiles").
		Columns(constants.CreateUserColumns...).
		Values(values...).
		PlaceholderFormat(sq.Dollar).
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
func (profileStore *ProfileStore) ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error) {
	sql, args, err := sq.Select(constants.ListProfilesColumns...).From("profiles").ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return []dto.ListProfiles{}, err
	}
	rows, err := profileStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing create project insert query:", err)
		return []dto.ListProfiles{}, err
	}

	for rows.Next() {
		var value dto.ListProfiles
		rows.Scan(&value.ID, &value.Name, &value.Email, &value.YearsOfExperience, &value.PrimarySkills, &value.IsCurrentEmployee)

		values = append(values, value)
	}
	defer rows.Close()

	return values, nil
}

// GetProfile returns a details profile in the Database that are currently available for perticular ID
func (profileStore *ProfileStore) GetProfile(ctx context.Context, profileID int) (value dto.ResponseProfile, err error) {
	query, args, err := sq.Select(constants.ResponseProfileColumns...).From("profiles").
		Where(sq.Eq{"id": profileID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return dto.ResponseProfile{}, err
	}

	rows, err := profileStore.db.Query(ctx, query, args...)
	if err != nil {
		zap.S().Error("Error executing get profile query: ", err)
		return dto.ResponseProfile{}, err
	}

	if rows.Next() {
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Email, &value.Gender, &value.Mobile, &value.Designation, &value.Description, &value.Title, &value.YearsOfExperience, &value.PrimarySkills, &value.SecondarySkills, &value.GithubLink, &value.LinkedinLink); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return dto.ResponseProfile{}, err
		}
	} else {
		zap.S().Info("No profile found for profileID: ", profileID)
		return dto.ResponseProfile{}, errors.ErrNoRecordFound
	}
	defer rows.Close()

	return value, nil
}

// UpdateProfile updates an existing user profile in the database.
func (profileStore *ProfileStore) UpdateProfile(ctx context.Context, profileID int, pd dto.UpdateProfileRequest) (int, error) {

	updateQuery, args, err := sq.Update("profiles").
		Set("name", pd.Profile.Name).Set("email", pd.Profile.Email).Set("gender", pd.Profile.Gender).
		Set("mobile", pd.Profile.Mobile).Set("designation", pd.Profile.Designation).Set("description", pd.Profile.Description).Set("title", pd.Profile.Title).Set("years_of_experience", pd.Profile.YearsOfExperience).Set("primary_skills", pd.Profile.PrimarySkills).Set("secondary_skills", pd.Profile.SecondarySkills).Set("github_link", pd.Profile.GithubLink).
		Set("linkedin_link", pd.Profile.LinkedinLink).Where(sq.Eq{"id": profileID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating profile update query: ", err)
		return 0, err
	}

	res, err := profileStore.db.Exec(ctx, updateQuery, args...)
	if err != nil {
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
