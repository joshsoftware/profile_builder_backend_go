package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
)

type ProfileStore struct{
	db *pgx.Conn
}

type ProfileStorer interface {
	CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error
}

func NewProfileRepo(db *pgx.Conn) ProfileStorer{
	return &ProfileStore{
        db: db,
    }
}

func (profileStore *ProfileStore) CreateProfile(pd dto.CreateProfileRequest, ctx context.Context) error{

	values := []interface{}{pd.Profile.Name, pd.Profile.Email, pd.Profile.Gender, pd.Profile.Mobile, pd.Profile.Designation, pd.Profile.Description, pd.Profile.Title, pd.Profile.YearsOfExperience, pd.Profile.PrimarySkills, pd.Profile.SecondarySkills, pd.Profile.GithubLink, pd.Profile.LinkedinLink, 1, 1}

	insertQuery, args ,err :=  sq.Insert("profiles").
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
