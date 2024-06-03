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
	CreateCertificate(ctx context.Context, values []CertificateRepo) error
	UpdateCertificate(ctx context.Context, profileID int, eduID int, req UpdateCertificateRepo) (int, error)
}

// NewCertificateRepo creates a new instance of CertificateRepo.
func NewCertificateRepo(db *pgx.Conn) CertificateStorer {
	return &CertificateStore{
		db: db,
	}
}

// CreateCertificate inserts certificate details into the database.
func (certificateStore *CertificateStore) CreateCertificate(ctx context.Context, values []CertificateRepo) error {

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
	_, err = certificateStore.db.Exec(ctx, insertQuery, args...)
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

// UpdateCertificate updates certificates details into the database.
func (certificateStore *CertificateStore) UpdateCertificate(ctx context.Context, profileID int, eduID int, req UpdateCertificateRepo) (int, error) {
	updateQuery, args, err := sq.Update("certificates").Set("name", req.Name).Set("organization_name", req.OrganizationName).Set("description", req.Description).Set("issued_date", req.IssuedDate).
		Set("from_date", req.FromDate).Set("to_date", req.ToDate).Set("updated_at", req.UpdatedAt).Set("updated_by_id", req.UpdatedByID).Where(sq.Eq{"profile_id": profileID, "id": eduID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating certificates update query: ", err)
		return 0, err
	}

	res, err := certificateStore.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing certificates update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : certificates")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}
