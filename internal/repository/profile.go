package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

type ProfileStore struct {
	db *pgx.Conn
}

type ProfileStorer interface {
	CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) error
	CreateEducation(ctx context.Context, values []EducationDao) error
	CreateProject(ctx context.Context, values []ProjectDao) error
}

func NewProfileRepo(db *pgx.Conn) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

func (profileStore *ProfileStore) CreateProfile(ctx context.Context, pd dto.CreateProfileRequest) error {

	values := []interface{}{pd.Profile.Name, pd.Profile.Email, pd.Profile.Gender, pd.Profile.Mobile, pd.Profile.Designation, pd.Profile.Description, pd.Profile.Title, pd.Profile.YearsOfExperience, pd.Profile.PrimarySkills, pd.Profile.SecondarySkills, pd.Profile.GithubLink, pd.Profile.LinkedinLink, 1, 1}

	insertQuery, args, err := sq.Insert("profiles").
		Columns(constants.CreateUserColumns...).
		Values(values...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		zap.S().Error("Error generating profile insert query: ", err)
		return err
	}

	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.New("profile already exists")
		}
		zap.S().Error("Error executing insert query:", err)
		return err
	}

	return nil
}

func (profileStore *ProfileStore) CreateEducation(ctx context.Context, values []EducationDao) error {

	for i := 0; i < len(values); i++ {

		value := []interface{}{
			values[i].Degree, values[i].UniversityName, values[i].Place, values[i].PercentageOrCgpa,
			values[i].PassingYear, values[i].CreatedAt, values[i].UpdatedAt, values[i].CreatedById,
			values[i].UpdatedById, values[i].ProfileId,
		}

		insertQuery, args, err := sq.Insert("educations").
			Columns(constants.CreateEducationColumns...).
			Values(value...).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			zap.S().Error("Error generating education insert query: ", err)
			return err
		}

		_, err = profileStore.db.Exec(ctx, insertQuery, args...)
		if err != nil {
			if helpers.IsDuplicateKeyError(err) {
				return errors.New("education already exists")
			}
			zap.S().Error("error executing create education insert query:", err)
			return err
		}
	}
	return nil
}

func (profileStore *ProfileStore) CreateProject(ctx context.Context, values []ProjectDao) error {

	for i := 0; i < len(values); i++ {

		value := []interface{}{
			values[i].Name, values[i].Description, values[i].Role, values[i].Responsibilities,
			values[i].Technologies, values[i].TechWorkedOn, values[i].WorkingStartDate, values[i].WorkingEndDate, values[i].Duration, values[i].CreatedAt, values[i].UpdatedAt, values[i].CreatedById, values[i].UpdatedById, values[i].ProfileId,
		}

		insertQuery, args, err := sq.Insert("projects").
			Columns(constants.CreateProjectColumns...).
			Values(value...).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			zap.S().Error("Error generating project insert query: ", err)
			return err
		}

		_, err = profileStore.db.Exec(ctx, insertQuery, args...)
		if err != nil {
			if helpers.IsDuplicateKeyError(err) {
				return errors.New("project already exists")
			}
			zap.S().Error("error executing create project insert query:", err)
			return err
		}
	}
	return nil
}
