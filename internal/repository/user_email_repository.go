package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

type EmailStore struct {
	db *pgxpool.Pool
}

type EmailStorer interface {
	GetEmailByProfileID(ctx context.Context, sendRequest SendUserInvitationRequest, tx pgx.Tx) (string, error)
	GetCreatedByIdByProfileID(ctx context.Context, sendRequest SendUserInvitationRequest, tx pgx.Tx) (int, error)
	GetUserEmailByUserID(ctx context.Context, userID int, tx pgx.Tx) (string, error)
	CreateSendInvitation(ctx context.Context, sendRequest EmailRepo, tx pgx.Tx) error
	UpdateProfileCompleteStatus(ctx context.Context, updateReq UpadateRequest, tx pgx.Tx) error
}

func NewUserEmailRepo(db *pgxpool.Pool) EmailStorer {
	return &EmailStore{
		db: db,
	}
}

func (userStore *EmailStore) GetEmailByProfileID(ctx context.Context, sendRequest SendUserInvitationRequest, tx pgx.Tx) (string, error) {
	queryBuilder := psql.Select("email").From("profiles").Where(sq.Eq{"id": sendRequest.ProfileID})
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return "", err
	}

	var email string
	err = tx.QueryRow(ctx, sql, args...).Scan(&email)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.S().Infof("No profile found with profile ID %d", sendRequest.ProfileID)
			return "", err
		}
		zap.S().Error("Error executing query: ", err)
		return "", err
	}
	return email, nil
}

func (userStore *EmailStore) GetCreatedByIdByProfileID(ctx context.Context, sendRequest SendUserInvitationRequest, tx pgx.Tx) (int, error) {
	queryBuilder := psql.Select("created_by_id").From("invitations").Where(sq.And{sq.Eq{"profile_id": sendRequest.ProfileID}, sq.Eq{"is_profile_complete": 0}})
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return 0, err
	}
	var createdByID int
	err = tx.QueryRow(ctx, sql, args...).Scan(&createdByID)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.S().Info("No invitations found with profile ID %d", sendRequest.ProfileID)
			return 0, err
		}
		zap.S().Error("Error executing query: ", err)
		return 0, err
	}
	return createdByID, nil
}

func (userStore *EmailStore) GetUserEmailByUserID(ctx context.Context, userID int, tx pgx.Tx) (string, error) {
	queryBuilder := psql.Select("email").From("users").Where(sq.Eq{"id": userID})
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return "", err
	}

	var email string
	err = tx.QueryRow(ctx, sql, args...).Scan(&email)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.S().Infof("No user found with user ID %d", userID)
			return "", err
		}
		zap.S().Error("Error executing query: ", err)
		return "", err
	}
	return email, nil
}

func (emailStore *EmailStore) CreateSendInvitation(ctx context.Context, emailRepo EmailRepo, tx pgx.Tx) error {
	values := []interface{}{
		emailRepo.ProfileID, emailRepo.ProfileComplete, emailRepo.CreatedAt, emailRepo.UpdatedAt,
		emailRepo.CreatedByID, emailRepo.UpdatedByID,
	}

	insertQuery, args, err := psql.Insert("invitations").
		Columns(constants.RequestInvitationColumns...).
		Values(values...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		zap.S().Error("Error generating email invitation insert query: ", err)
		return err
	}

	if tx == nil {
		zap.S().Error("Transaction is nil")
		return fmt.Errorf("internal error: transaction is nil")
	}

	res, err := tx.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) || helpers.IsInvalidProfileError(err) {
			zap.S().Info("Error creating send invitation: ", err)
			return err
		}
		zap.S().Error("Error executing insert query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Info("No rows affected with given query : ", insertQuery)
		return nil
	}

	return nil
}

func (emailStore *EmailStore) UpdateProfileCompleteStatus(ctx context.Context, updateReq UpadateRequest, tx pgx.Tx) error {
	queryBuilder := psql.Update("invitations").Set("is_profile_complete", 1).Set("updated_at", updateReq.UpdatedAt).Where(sq.Eq{"profile_id": updateReq.ProfileID, "is_profile_complete": 0})
	updateQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating update query: ", err)
		return err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		zap.S().Error("Error executing update query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Info("No rows were updated for profile_id: ", updateReq.ProfileID)
		return nil
	}
	return nil
}
