package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// CertificateService represents a set of methods for accessing the certificates.
type CertificateService interface {
	CreateCertificate(ctx context.Context, cDetail specs.CreateCertificateRequest, profileID int, userID int) (ID int, err error)
	UpdateCertificate(ctx context.Context, profileID int, certID int, userID int, req specs.UpdateCertificateRequest) (ID int, err error)
	ListCertificates(ctx context.Context, profileID int, fitler specs.ListCertificateFilter) (value []specs.CertificateResponse, err error)
}

// CreateCerticate : Service layer function adds certicates details to a user profile.
func (certificateSvc *service) CreateCertificate(ctx context.Context, cDetail specs.CreateCertificateRequest, profileID int, userID int) (ID int, err error) {
	tx, _ := certificateSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := certificateSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	if len(cDetail.Certificates) == 0 {
		zap.S().Error("certificates payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.CertificateRepo
	for _, certificate := range cDetail.Certificates {
		val := repository.CertificateRepo{
			ProfileID:        profileID,
			Name:             certificate.Name,
			OrganizationName: certificate.OrganizationName,
			Description:      certificate.Description,
			IssuedDate:       certificate.IssuedDate,
			FromDate:         certificate.FromDate,
			ToDate:           certificate.ToDate,
			CreatedAt:        today,
			UpdatedAt:        today,
			CreatedByID:      userID,
			UpdatedByID:      userID,
		}

		value = append(value, val)
	}

	err = certificateSvc.CertificateRepo.CreateCertificate(ctx, value, tx)
	if err != nil {
		zap.S().Error("Unable to create Certificate : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return profileID, nil
}

// UpdateCertificate in the service layer update a certificates of specific profile.
func (certificateSvc *service) UpdateCertificate(ctx context.Context, profileID int, certID int, userID int, req specs.UpdateCertificateRequest) (ID int, err error) {
	tx, _ := certificateSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := certificateSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	var value repository.UpdateCertificateRepo
	value.Name = req.Certificate.Name
	value.OrganizationName = req.Certificate.OrganizationName
	value.Description = req.Certificate.Description
	value.IssuedDate = req.Certificate.IssuedDate
	value.FromDate = req.Certificate.FromDate
	value.ToDate = req.Certificate.ToDate
	value.UpdatedAt = today
	value.UpdatedByID = userID

	profileID, err = certificateSvc.CertificateRepo.UpdateCertificate(ctx, profileID, certID, value, tx)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("certificate(s) update with profile id : ", profileID)

	return profileID, nil
}

func (certificateSvc *service) ListCertificates(ctx context.Context, profileID int, filter specs.ListCertificateFilter) (value []specs.CertificateResponse, err error) {
	tx, _ := certificateSvc.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := certificateSvc.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	value, err = certificateSvc.CertificateRepo.ListCertificates(ctx, profileID, filter, tx)
	if err != nil {
		zap.S().Error("error to get certificate : ", err, " for profile id : ", profileID)
		return []specs.CertificateResponse{}, err
	}

	return value, nil
}
