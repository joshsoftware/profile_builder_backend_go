package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// BaseRepository implements the Trasanctions interface.
type BaseRepository struct {
	DB *pgx.Conn
}

// Trasanctions defines methods to perform transactions.
type Trasanctions interface {
	BeginTx(ctx context.Context) (tx pgx.Tx, err error)
	CommitTx(tx pgx.Tx) (err error)
	RollbackTx(tx pgx.Tx) (err error)
	HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error)
}

// BeginTx used to Begin specific transaction
func (repo *BaseRepository) BeginTx(ctx context.Context) (tx pgx.Tx, err error) {

	sqlDB, err := repo.DB.Begin(ctx)
	if err != nil {
		zap.S().Errorf("error occured while initiating database transaction: %v", err.Error())
		return nil, err
	}

	return sqlDB, nil
}

// RollbackTx used to rollback specific transaction
func (repo *BaseRepository) RollbackTx(ctx context.Context, tx pgx.Tx) (err error) {
	err = tx.Rollback(ctx)
	return
}

// CommitTx used to commit specific transaction
func (repo *BaseRepository) CommitTx(ctx context.Context, tx pgx.Tx) (err error) {
	err = tx.Commit(ctx)
	return
}

// HandleTransaction used to handle transaction
func (repo *BaseRepository) HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error) {
	if incomingErr != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return
		}
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		return
	}
	return
}
