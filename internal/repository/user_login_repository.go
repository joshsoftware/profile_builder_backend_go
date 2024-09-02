package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"

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
	GetUserInfo(ctx context.Context, filter specs.UserInfoFilter) (User, error)
	CreateUser(ctx context.Context, name, email string, role string, tx pgx.Tx) error
	RemoveUser(ctx context.Context, email string, tx pgx.Tx) error
}

// NewUserLoginRepo defines repo dependancies
func NewUserLoginRepo(db *pgxpool.Pool) UserStorer {
	return &UserStore{
		db: db,
	}
}

// GetUserInfo returns and checks the id and role of a specific email or id
func (userStore *UserStore) GetUserInfo(ctx context.Context, filter specs.UserInfoFilter) (User, error) {
	var user User

	selectBuilder := sq.Select(constants.RequestUserColumns...).From(userTable).PlaceholderFormat(sq.Dollar)
	if filter.Email != "" {
		selectBuilder = selectBuilder.Where(sq.Eq{"email": filter.Email})
	} else if filter.ID > 0 {
		selectBuilder = selectBuilder.Where(sq.Eq{"id": filter.ID})
	} else {
		return User{}, errors.New("filter must contain either email or id")
	}

	selectQuery, args, err := selectBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return User{}, err
	}

	// execute the query using pgx
	row := userStore.db.QueryRow(ctx, selectQuery, args...)
	err = row.Scan(&user.ID, &user.Email, &user.Role, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, errs.ErrNoRecordFound
		}
		if helpers.IsInvalidProfileError(err) {
			zap.S().Error("Invalid profile error : ", err)
			return User{}, errs.ErrInvalidProfile
		}

		zap.S().Error("Error executing select query : ", err)
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user with the given email and role
func (userStore *UserStore) CreateUser(ctx context.Context, name, email string, role string, tx pgx.Tx) error {
	checkQuery := psql.Select("1").From(userTable).Where(sq.Eq{"email": email}).Limit(1)
	checkSQL, checkArgs, err := checkQuery.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return err
	}

	var exists int64
	err = tx.QueryRow(ctx, checkSQL, checkArgs...).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		zap.S().Error("Error executing select query: ", err)
		return err
	}

	if exists > 0 {
		zap.S().Info("Record already exists for email: ", email)
		return nil
	}

	query := psql.Insert(userTable).Columns("email", "role", "name").Values(email, role, name).Suffix("RETURNING id")
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
			return errs.ErrDuplicateKey
		} else if helpers.IsInvalidProfileError(err) {
			zap.S().Info("Error creating user: ", err)
			return err
		}
		zap.S().Error("Error executing insert query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Info("No rows affected with given query: ", sql)
		return errs.ErrNoRecordFound
	}
	return nil
}

// RemoveUser removes a user with the given email
func (userStore *UserStore) RemoveUser(ctx context.Context, email string, tx pgx.Tx) error {
	var role string
	roleQuery := psql.Select("role").From(userTable).Where(sq.Eq{"email": email})
	sql, args, err := roleQuery.ToSql()
	if err != nil {
		zap.S().Error("Error generating role query: ", err)
		return err
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&role)
	if err != nil {
		zap.S().Error("Error executing role query: ", err)
		return err
	}

	if role == constants.Admin {
		zap.S().Info("Cannot remove user with admin role")
		return nil
	}

	query := psql.Delete(userTable).Where(sq.Eq{"email": email})
	sql, args, err = query.ToSql()
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
		zap.S().Info("No rows affected for remove user")
		return errs.ErrNoRecordFound
	}
	return nil
}
