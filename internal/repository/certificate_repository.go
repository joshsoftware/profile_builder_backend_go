package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// CertificateStore implements the CertificateStorer interface.
type CertificateStore struct {
	db *pgx.Conn
}

// CertificateStorer defines methods to interact with user certificate ralated data.
type CertificateStorer interface {
	CreateCertificate(ctx context.Context, values []CertificateDao) error
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewCertificateRepo(db *pgx.Conn) CertificateStorer {
	return &CertificateStore{
		db: db,
	}
}

// CreateCertificate inserts certificate details into the database.
func (profileStore *CertificateStore) CreateCertificate(ctx context.Context, values []CertificateDao) error {

	insertBuilder := sq.Insert("certificates").
		Columns(constants.CreateCertificateColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, value := range values {
		insertBuilder = insertBuilder.Values(
			value.Name, value.OrganizationName, value.Description, value.IssuedDate,
			value.FromDate, value.ToDate, value.CreatedAt, value.UpdatedAt,
			value.CreatedByID, value.UpdatedByID, value.ProfileID,
		)
	}

	insertQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating certificate insert query: ", err)
		return err
	}
	_, err = profileStore.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing create certificate insert query: ", err)
		return err
	}

	return nil
}
