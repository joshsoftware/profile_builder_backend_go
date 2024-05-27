package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type CertificateService interface {
	CreateCertificate(ctx context.Context, cDetail dto.CreateCertificateRequest) (profileID int, err error)
}

// CreateCerticate : Service layer function adds certicates details to a user profile.
func (profileSvc *service) CreateCertificate(ctx context.Context, cDetail dto.CreateCertificateRequest) (profileID int, err error) {
	if len(cDetail.Certificates) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.CertificateDao
	for i := 0; i < len(cDetail.Certificates); i++ {
		var val repository.CertificateDao

		val.ProfileID = cDetail.ProfileID
		val.Name = cDetail.Certificates[i].Name
		val.OrganizationName = cDetail.Certificates[i].OrganizationName
		val.Description = cDetail.Certificates[i].Description
		val.IssuedDate = cDetail.Certificates[i].IssuedDate
		val.FromDate = cDetail.Certificates[i].FromDate
		val.ToDate = cDetail.Certificates[i].ToDate
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.CertificateRepo.CreateCertificate(ctx, value)
	if err != nil {
		return 0, err
	}

	return cDetail.ProfileID, nil
}
