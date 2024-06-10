package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
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
	ListCertificates(ctx context.Context, profileID int) ([]dto.CertificateResponse, error)
}

// NewCertificateRepo creates a new instance of CertificateRepo.
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

// GetCertificates fetches certificates details from the database.
func (certificateStore *CertificateStore) ListCertificates(ctx context.Context, profileID int) (values []dto.CertificateResponse, err error) {
	sql, args, err := sq.Select(constants.ResponseCertificatesColumns...).From("certificates").Where(sq.Eq{"profile_id": profileID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating get certificates query: ", err)
		return []dto.CertificateResponse{}, err
	}

	rows, err := certificateStore.db.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get certificates query: ", err)
		return []dto.CertificateResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var val dto.CertificateResponse
		err = rows.Scan(&val.ProfileID, &val.Name, &val.OrganizationName, &val.Description, &val.IssuedDate, &val.FromDate, &val.ToDate)
		if err != nil {
			zap.S().Error("Error scanning certificates rows: ", err)
			return []dto.CertificateResponse{}, err
		}
		values = append(values, val)
	}

	if len(values) == 0 {
		zap.S().Error("No certificates found for profile id: ", profileID)
		return []dto.CertificateResponse{}, errors.ErrNoRecordFound
	}

	return values, nil
}
