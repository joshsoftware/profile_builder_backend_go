package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// struct to store the database connection
type UserStore struct {
	db *pgx.Conn
}

// interface
type UserStorer interface {
	GetUserIdByEmail(ctx context.Context, email string) (int64, error)
}

func NewUserLoginRepo(db *pgx.Conn) UserStorer {
	return &UserStore{
		db: db,
	}
}

func (profileStore *UserStore) GetUserIdByEmail(ctx context.Context, email string) (int64, error) {

	var user UserDao

	// used squirrel to build the query
	selectBuilder := sq.Select("id").From("users").Where(sq.Eq{"email": email}).PlaceholderFormat(sq.Dollar)

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
		if err == sql.ErrNoRows {
			return 0, errors.ErrNoRecordFound
		} else {
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
	}

	return user.ID, nil

}
