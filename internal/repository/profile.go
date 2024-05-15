package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
)

type ProfileStore struct {
	db *pgx.Conn
}

type ProfileStorer interface {
	CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error
	GetProfileByEmail(ctx context.Context, email string) (int64, error)
}

func NewProfileRepo(db *pgx.Conn) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

func (profileStore *ProfileStore) CreateProfile(pd dto.CreateProfileRequest, ctx context.Context) error {

	values := []interface{}{pd.Profile.Name, pd.Profile.Email, pd.Profile.Gender, pd.Profile.Mobile, pd.Profile.Designation, pd.Profile.Description, pd.Profile.Title, pd.Profile.YearsOfExperience, pd.Profile.PrimarySkills, pd.Profile.SecondarySkills, pd.Profile.GithubLink, pd.Profile.LinkedinLink, 1, 1}

	insertQuery, args, err := sq.Insert("profiles").
		Columns(constants.CreateUserColumns...).
		Values(values...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		fmt.Println("Error generating profile insert query: ", err)
		return err
	}

	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.New("profile already exists")
		}
		fmt.Println("Error executing insert query:", err)
		return err
	}

	return nil
}

func (profileStore *ProfileStore) GetProfileByEmail(ctx context.Context, email string) (int64, error) {
	var user dto.User
	row := profileStore.db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no user found with this email")
		} else {
			return 0, err
		}
	}
	return user.ID, nil
}
