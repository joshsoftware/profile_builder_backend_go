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
	CreateEducation(ctx context.Context, values []EducationDao) error
	CreateProject(ctx context.Context, values []ProjectDao) error
	CreateExperience(ctx context.Context, values []ExperienceDao) error
	CreateCertificate(ctx context.Context, values []CertificateDao) error
	CreateAchievement(ctx context.Context, values []AchievementDao) error

	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
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

// CreateEducation inserts education details into the database.
func (profileStore *ProfileStore) CreateEducation(ctx context.Context, values []EducationDao) error {

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

// CreateProject inserts project details into the database.
func (profileStore *ProfileStore) CreateProject(ctx context.Context, values []ProjectDao) error {

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

// CreateExperience inserts experience details into the database.
func (profileStore *ProfileStore) CreateExperience(ctx context.Context, values []ExperienceDao) error {

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

// CreateCertificate inserts certificate details into the database.
func (profileStore *ProfileStore) CreateCertificate(ctx context.Context, values []CertificateDao) error {

	insertBuilder := sq.Insert("certificates").
		Columns(constants.CreateCertificateColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.OrganizationName, value.Description, value.IssuedDate,
			value.FromDate, value.ToDate, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating certificate insert query: ", err)
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
		zap.S().Error("Error executing create certificate insert query: ", err)
		return err
	}

	return nil
}

// CreateAchievement inserts achievements details into the database.
func (profileStore *ProfileStore) CreateAchievement(ctx context.Context, values []AchievementDao) error {

	insertBuilder := sq.Insert("achievements").
		Columns(constants.CreateAchievementColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.Description, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating achievement insert query: ", err)
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
		zap.S().Error("Error executing create achievement insert query: ", err)
		return err
	}

	return nil
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
