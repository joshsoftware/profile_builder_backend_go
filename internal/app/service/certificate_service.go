package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// CertificateService represents a set of methods for accessing the certificates.
type CertificateService interface {
	CreateCertificate(ctx context.Context, cDetail specs.CreateCertificateRequest, ID int) (profileID int, err error)
	UpdateCertificate(ctx context.Context, profileID string, eduID string, req specs.UpdateCertificateRequest) (id int, err error)
	ListCertificates(ctx context.Context, profileID int, fitler specs.ListCertificateFilter) (value []specs.CertificateResponse, err error)
}

// CreateCerticate : Service layer function adds certicates details to a user profile.
func (certificateSvc *service) CreateCertificate(ctx context.Context, cDetail specs.CreateCertificateRequest, ID int) (profileID int, err error) {
	if len(cDetail.Certificates) == 0 {
		zap.S().Error("certificates payload can't be empty")
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.CertificateRepo
	for i := 0; i < len(cDetail.Certificates); i++ {
		var val repository.CertificateRepo

		val.ProfileID = ID
		val.Name = cDetail.Certificates[i].Name
		val.OrganizationName = cDetail.Certificates[i].OrganizationName
		val.Description = cDetail.Certificates[i].Description
		val.IssuedDate = cDetail.Certificates[i].IssuedDate
		val.FromDate = cDetail.Certificates[i].FromDate
		val.ToDate = cDetail.Certificates[i].ToDate
		val.CreatedAt = today
		val.UpdatedAt = today
		//TODO by context
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = certificateSvc.CertificateRepo.CreateCertificate(ctx, value)
	if err != nil {
		zap.S().Error("Unable to create Certificate : ", err, " for profile id : ", profileID)
		return 0, err
	}

	return ID, nil
}

// UpdateCertificate in the service layer update a certificates of specific profile.
func (certificateSvc *service) UpdateCertificate(ctx context.Context, profileID string, eduID string, req specs.UpdateCertificateRequest) (id int, err error) {
	pid, id, err := helpers.MultipleConvertStringToInt(profileID, eduID)
	if err != nil {
		zap.S().Error("error to get certificate params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	var value repository.UpdateCertificateRepo
	value.Name = req.Certificate.Name
	value.OrganizationName = req.Certificate.OrganizationName
	value.Description = req.Certificate.Description
	value.IssuedDate = req.Certificate.IssuedDate
	value.FromDate = req.Certificate.FromDate
	value.ToDate = req.Certificate.ToDate
	value.UpdatedAt = today
	//TODO by context
	value.UpdatedByID = 1

	id, err = certificateSvc.CertificateRepo.UpdateCertificate(ctx, pid, id, value)
	if err != nil {
		zap.S().Error("Unable to update education : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("certificate(s) update with profile id : ", id)

	return id, nil
}

func (certificateSvc *service) ListCertificates(ctx context.Context, profileID int, filter specs.ListCertificateFilter) (value []specs.CertificateResponse, err error) {
	value, err = certificateSvc.CertificateRepo.ListCertificates(ctx, profileID, filter)
	if err != nil {
		zap.S().Error("error to get certificate : ", err, " for profile id : ", profileID)
		return []specs.CertificateResponse{}, err
	}

	return value, nil
}
