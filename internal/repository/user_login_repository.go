package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// UserStore implements the UserStorer interface.
type UserStore struct {
	db *pgxpool.Pool
}

var (
	userTable = "users"
)

// UserStorer defines methods to interact with user data.
type UserStorer interface {
	GetUserIDByEmail(ctx context.Context, email string) (int64, error)
}

// NewUserLoginRepo defines repo dependancies
func NewUserLoginRepo(db *pgxpool.Pool) UserStorer {
	return &UserStore{
		db: db,
	}
}

// GetUserIDByEmail returns and checks the id of specific email
func (profileStore *UserStore) GetUserIDByEmail(ctx context.Context, email string) (int64, error) {

	var user UserDao

	// used squirrel to build the query
	selectBuilder := sq.Select("id").From(userTable).Where(sq.Eq{"email": email}).PlaceholderFormat(sq.Dollar)

	// generate the query
	selectQuery, args, err := selectBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return 0, err
	}

	// execute the query using pgx
	row := profileStore.db.QueryRow(ctx, selectQuery, args...)
	err = row.Scan(&user.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, errors.ErrNoRecordFound
		}
		if helpers.IsDuplicateKeyError(err) {
			zap.S().Error("Duplicate key error : ", err)
			return 0, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			zap.S().Error("Invalid profile error : ", err)
			return 0, errors.ErrInvalidProfile
		}

		zap.S().Error("Error executing select query : ", err)
	}

	return user.ID, nil
}
