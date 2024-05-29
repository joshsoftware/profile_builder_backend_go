package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateCertificate(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	var repodeps = service.RepoDeps{
		CertificateDeps: mockCertificateRepo,
	}
	certificateService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           dto.CreateCertificateRequest
		setup           func(certificateMock *mocks.CertificateStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for certificate details",
			input: dto.CreateCertificateRequest{
				ProfileID: 1,
				Certificates: []dto.Certificate{
					{
						Name:             "Full Stack GO Certificate",
						OrganizationName: "Josh Software",
						Description:      "Description about certificate",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed because CreateCertificate returns an error",
			input: dto.CreateCertificateRequest{
				ProfileID: 10000,
				Certificates: []dto.Certificate{
					{
						Name:             "Full Stack Data Science Certificate",
						OrganizationName: "Balambika Team",
						Description:      "Description",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of missing certificate name",
			input: dto.CreateCertificateRequest{
				ProfileID: 1,
				Certificates: []dto.Certificate{
					{
						Name:             "",
						OrganizationName: "Org C",
						Description:      "Description C",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateDao")).Return(errors.New("Missing certificate name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed because of empty payload",
			input: dto.CreateCertificateRequest{
				ProfileID:    1,
				Certificates: []dto.Certificate{},
			},
			setup:           func(certificateMock *mocks.CertificateStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateRepo)

			// Test the service
			_, err := certificateService.CreateCertificate(context.TODO(), test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
