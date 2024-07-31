package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
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
	GetUserInfoByEmail(ctx context.Context, email string) (UserDao, error)
	CreateUserAsEmployee(ctx context.Context, email string, tx pgx.Tx) error
	RemoveUserEmployee(ctx context.Context, email string, tx pgx.Tx) error
}

// NewUserLoginRepo defines repo dependancies
func NewUserLoginRepo(db *pgxpool.Pool) UserStorer {
	return &UserStore{
		db: db,
	}
}

// GetUserIDByEmail returns and checks the id of specific email
func (profileStore *UserStore) GetUserInfoByEmail(ctx context.Context, email string) (UserDao, error) {

	var user UserDao

	// used squirrel to build the query
	selectBuilder := sq.Select("id", "role").From(userTable).Where(sq.Eq{"email": email}).PlaceholderFormat(sq.Dollar)

	// generate the query
	selectQuery, args, err := selectBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return UserDao{}, err
	}

	// execute the query using pgx
	row := profileStore.db.QueryRow(ctx, selectQuery, args...)
	err = row.Scan(&user.ID, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return UserDao{}, errors.ErrNoRecordFound
		}
		if helpers.IsDuplicateKeyError(err) {
			zap.S().Error("Duplicate key error : ", err)
			return UserDao{}, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			zap.S().Error("Invalid profile error : ", err)
			return UserDao{}, errors.ErrInvalidProfile
		}

		zap.S().Error("Error executing select query : ", err)
	}

	return user, nil
}

// CreateUserAsEmployee creates a new user with the given email and role as employee
func (userStore *UserStore) CreateUserAsEmployee(ctx context.Context, email string, tx pgx.Tx) error {
	query := psql.Insert(userTable).Columns("email", "role").Values(email, constants.Employee).Suffix("RETURNING id")
	sql, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("Error generating insert query: ", err)
		return err
	}

	if tx == nil {
		zap.S().Error("Transaction is nil")
		return fmt.Errorf("internal error: transaction is nil")
	}

	res, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			zap.S().Info("Record already exists for email: ", email)
			return errors.ErrDuplicateKey
		} else if helpers.IsInvalidProfileError(err) {
			zap.S().Info("Error creating user: ", err)
			return err
		}
		zap.S().Error("Error executing insert query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Info("No rows affected with given query: ", sql)
		return errors.ErrNoRecordFound
	}
	return nil
}

func (userStore *UserStore) RemoveUserEmployee(ctx context.Context, email string, tx pgx.Tx) error {
	query := psql.Delete(userTable).Where(sq.Eq{"email": email})
	sql, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("Error generating delete query: ", err)
		return err
	}

	if tx == nil {
		zap.S().Error("Transaction is nil")
		return fmt.Errorf("internal error: transaction is nil")
	}

	res, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing delete query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Info("No rows affected with given query: ", sql)
		return errors.ErrNoRecordFound
	}
	return nil
}
