package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

type EmailStore struct {
	db *pgxpool.Pool
}

type EmailStorer interface {
	GetInvitations(ctx context.Context, profileID int, tx pgx.Tx) (specs.InvitationResponse, error)
	CreateInvitation(ctx context.Context, invitation Invitations, tx pgx.Tx) error
	UpdateProfileCompleteStatus(ctx context.Context, profileID int, updateReq UpadateRequest, tx pgx.Tx) error
}

func NewUserEmailRepo(db *pgxpool.Pool) EmailStorer {
	return &EmailStore{
		db: db,
	}
}

func (userStore *EmailStore) GetInvitations(ctx context.Context, profileID int, tx pgx.Tx) (specs.InvitationResponse, error) {
	queryBuilder := psql.Select(constants.RequestInvitationColumns...).From("invitations").Where(sq.And{sq.Eq{"profile_id": profileID}, sq.Eq{"is_profile_complete": 0}})
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query for get invitations: ", err)
		return specs.InvitationResponse{}, err
	}
	var Invitation specs.InvitationResponse
	err = tx.QueryRow(ctx, sql, args...).Scan(
		&Invitation.ProfileID, &Invitation.ProfileComplete, &Invitation.CreatedAt, &Invitation.UpdatedAt, &Invitation.CreatedByID, &Invitation.UpdatedByID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.S().Info("No invitations found with profile ID %d", profileID)
			return specs.InvitationResponse{}, err
		}
		zap.S().Error("Error executing query for get invitations: ", err)
		return specs.InvitationResponse{}, err
	}
	return Invitation, nil
}

func (emailStore *EmailStore) CreateInvitation(ctx context.Context, invitation Invitations, tx pgx.Tx) error {
	values := []interface{}{
		invitation.ProfileID, invitation.ProfileComplete, invitation.CreatedAt, invitation.UpdatedAt,
		invitation.CreatedByID, invitation.UpdatedByID,
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

func (emailStore *EmailStore) UpdateProfileCompleteStatus(ctx context.Context, profileID int, updateReq UpadateRequest, tx pgx.Tx) error {
	queryBuilder := psql.Update("invitations").Set("is_profile_complete", updateReq.ProfileComplete).Set("updated_at", updateReq.UpdatedAt).Where(sq.Eq{"profile_id": profileID, "is_profile_complete": 0})
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
		zap.S().Info("No rows were updated for profile_id: ", profileID)
		return nil
	}
	return nil
}
