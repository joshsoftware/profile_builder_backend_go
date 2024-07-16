package specs

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// ListCertificateFilter used to filter certificates based on ids and names
type ListCertificateFilter struct {
	CertificateIDs []int    `json:"certificate_ids"`
	Names          []string `json:"names"`
}

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

// UpdateCertificateRequest struct represents a request to update a certificate
type UpdateCertificateRequest struct {
	Certificate Certificate `json:"certificate"`
}

// Validate func checks if the CreateCertificateRequest is valid.
func (req *CreateCertificateRequest) Validate() error {

	if len(req.Certificates) == 0 {
		return fmt.Errorf("%s : certificates ", errors.ErrEmptyPayload.Error())
	}

	for _, cert := range req.Certificates {
		if cert.Name == "" {
			return fmt.Errorf("%s : certificate name", errors.ErrParameterMissing.Error())
		}

		if cert.IssuedDate == "" {
			return fmt.Errorf("%s : issued date", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}

// Validate func checks if the UpdateCertificateRequest is valid.
func (req *UpdateCertificateRequest) Validate() error {

	fields := map[string]string{
		"name":              req.Certificate.Name,
		"organization name": req.Certificate.OrganizationName,
		"decsription":       req.Certificate.Description,
		"issued date":       req.Certificate.IssuedDate,
		"from date":         req.Certificate.FromDate,
		"to date":           req.Certificate.ToDate,
	}

	for fieldName, fieldValue := range fields {
		if fieldValue == "" {
			return fmt.Errorf("%s : %s", errors.ErrParameterMissing.Error(), fieldName)
		}
	}

	return nil
}

type DeleteCertificateRequest struct {
	ProfileID     int `json:"profile_id"`
	CertificateID int `json:"id"`
}
