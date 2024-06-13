package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateCertificateRequest struct represents a request to create certificates details.
type CreateCertificateRequest struct {
	ProfileID    int           `json:"profile_id"`
	Certificates []Certificate `json:"certificates"`
}

// Certificate struct represents details of an certificates.
type Certificate struct {
	Name             string `json:"name"`
	OrganizationName string `json:"organization_name"`
	Description      string `json:"description"`
	IssuedDate       string `json:"issued_date"`
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
}

// CertificateResponse struct represents details of an certificates for specific id.
type CertificateResponse struct {
	ID               int    `json:"id"`
	ProfileID        int    `json:"profile_id"`
	Name             string `json:"name"`
	OrganizationName string `json:"organization_name"`
	Description      string `json:"description"`
	IssuedDate       string `json:"issued_date"`
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
}

// ResponseCertificate used for response of certificates of profiles
type ResponseCertificate struct {
	Certificates []CertificateResponse `json:"certificates"`
}

// Validate func checks if the CreateCertificateRequest is valid.
func (req *CreateCertificateRequest) Validate() error {
	if req.ProfileID == 0 {
		return fmt.Errorf("%s : profile id", errors.ErrParameterMissing.Error())
	}
	if len(req.Certificates) == 0 {
		return fmt.Errorf("%s : certificates ", errors.ErrEmptyPayload.Error())
	}

	for _, edu := range req.Certificates {
		if edu.Name == "" {
			return fmt.Errorf("%s : certificate name", errors.ErrParameterMissing.Error())
		}

		if edu.OrganizationName == "" {
			return fmt.Errorf("%s : organization name", errors.ErrParameterMissing.Error())
		}

		if edu.Description == "" {
			return fmt.Errorf("%s : decsription", errors.ErrParameterMissing.Error())
		}

		if edu.IssuedDate == "" {
			return fmt.Errorf("%s : issued date", errors.ErrParameterMissing.Error())
		}

		if edu.FromDate == "" {
			return fmt.Errorf("%s : from date", errors.ErrParameterMissing.Error())
		}

		if edu.ToDate == "" {
			return fmt.Errorf("%s : to date", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}
