package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// CertificateStore implements the CertificateStorer interface.
type CertificateStore struct {
	db *pgxpool.Pool
}

// CertificateStorer defines methods to interact with user certificate ralated data.
type CertificateStorer interface {
	CreateCertificate(ctx context.Context, values []CertificateRepo, tx pgx.Tx) error
	UpdateCertificate(ctx context.Context, profileID int, eduID int, req UpdateCertificateRepo, tx pgx.Tx) (int, error)
	ListCertificates(ctx context.Context, profileID int, filter specs.ListCertificateFilter, tx pgx.Tx) ([]specs.CertificateResponse, error)
}

// NewCertificateRepo creates a new instance of CertificateRepo.
func NewCertificateRepo(db *pgxpool.Pool) CertificateStorer {
	return &CertificateStore{
		db: db,
	}
}

// CreateCertificate inserts certificate details into the database.
func (certificateStore *CertificateStore) CreateCertificate(ctx context.Context, values []CertificateRepo, tx pgx.Tx) error {
	insertBuilder := psql.Insert("certificates").
		Columns(constants.CreateCertificateColumns...)

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
	_, err = tx.Exec(ctx, insertQuery, args...)
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

// ListCertificates fetches certificates details from the database.
func (certificateStore *CertificateStore) ListCertificates(ctx context.Context, profileID int, filter specs.ListCertificateFilter, tx pgx.Tx) (values []specs.CertificateResponse, err error) {
	queryBuilder := sq.Select(constants.ResponseCertificatesColumns...).From("certificates").Where(sq.Eq{"profile_id": profileID})
	if len(filter.CertificateIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": filter.CertificateIDs})
	}
	if len(filter.Names) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"name": filter.Names})
	}
	sql, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		zap.S().Error("Error generating get certificates query: ", err)
		return []specs.CertificateResponse{}, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing get certificates query: ", err)
		return []specs.CertificateResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var val specs.CertificateResponse
		err = rows.Scan(&val.ID, &val.ProfileID, &val.Name, &val.OrganizationName, &val.Description, &val.IssuedDate, &val.FromDate, &val.ToDate)
		if err != nil {
			zap.S().Error("Error scanning certificates rows: ", err)
			return []specs.CertificateResponse{}, err
		}
		values = append(values, val)
	}

	return values, nil
}

// UpdateCertificate updates certificates details into the database.
func (certificateStore *CertificateStore) UpdateCertificate(ctx context.Context, profileID int, eduID int, req UpdateCertificateRepo, tx pgx.Tx) (int, error) {
	updateQuery, args, err := psql.Update("certificates").
		SetMap(map[string]interface{}{
			"name": req.Name, "organization_name": req.OrganizationName,
			"description": req.Description, "issued_date": req.IssuedDate,
			"from_date": req.FromDate, "to_date": req.ToDate,
			"updated_at": req.UpdatedAt, "updated_by_id": req.UpdatedByID,
		}).Where(sq.Eq{"profile_id": profileID, "id": eduID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating certificates update query: ", err)
		return 0, err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
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
