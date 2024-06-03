package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateCertificateRequest struct represents a request to create certificates details.
type CreateCertificateRequest struct {
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

// UpdateCertificateRequest struct represents a request to update a certificate
type UpdateCertificateRequest struct {
	Certificate Certificate `json:"certificate"`
}

// Validate func checks if the CreateCertificateRequest is valid.
func (req *CreateCertificateRequest) Validate() error {

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

// Validate func checks if the UpdateCertificateRequest is valid.
func (req *UpdateCertificateRequest) Validate() error {

	if req.Certificate.Name == "" {
		return fmt.Errorf("%s : certificate name", errors.ErrParameterMissing.Error())
	}

	if req.Certificate.OrganizationName == "" {
		return fmt.Errorf("%s : organization name", errors.ErrParameterMissing.Error())
	}

	if req.Certificate.Description == "" {
		return fmt.Errorf("%s : decsription", errors.ErrParameterMissing.Error())
	}

	if req.Certificate.IssuedDate == "" {
		return fmt.Errorf("%s : issued date", errors.ErrParameterMissing.Error())
	}

	if req.Certificate.FromDate == "" {
		return fmt.Errorf("%s : from date", errors.ErrParameterMissing.Error())
	}

	if req.Certificate.ToDate == "" {
		return fmt.Errorf("%s : to date", errors.ErrParameterMissing.Error())
	}

	return nil
}
