package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
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
		input           specs.CreateCertificateRequest
		setup           func(certificateMock *mocks.CertificateStorer)
		isErrorExpected bool
	}{
		{
			name: "Success_for_certificate_details",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{
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
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo")).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "Failed_because_createcertificate_returns_an_error",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{
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
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo")).Return(errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_missing_certificate_name",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{
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
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo")).Return(errors.New("Missing certificate name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_empty_payload",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{},
			},
			setup:           func(certificateMock *mocks.CertificateStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateRepo)

			// Test the service
			_, err := certificateService.CreateCertificate(context.TODO(), test.input, 1)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestUpdateCertificate(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	var repodeps = service.RepoDeps{
		CertificateDeps: mockCertificateRepo,
	}
	certService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       string
		certificateID   string
		input           specs.UpdateCertificateRequest
		setup           func(certificateMock *mocks.CertificateStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_updating_certificate_details",
			profileID:     "1",
			certificateID: "1",
			input: specs.UpdateCertificateRequest{
				Certificate: specs.Certificate{
					Name:             "Updated Certificate Name",
					OrganizationName: "Updated Organization",
					Description:      "Updated Description",
					IssuedDate:       "2023-01-01",
					FromDate:         "2022-01-01",
					ToDate:           "2024-01-01",
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo")).Return(1, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed_because_updatecertificate_returns_an_error",
			profileID:     "100000000000000000",
			certificateID: "1",
			input: specs.UpdateCertificateRequest{
				Certificate: specs.Certificate{
					Name:             "Certificate B",
					OrganizationName: "Organization B",
					Description:      "Description B",
					IssuedDate:       "2023-01-01",
					FromDate:         "2022-01-01",
					ToDate:           "2024-01-01",
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("UpdateCertificate", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("repository.UpdateCertificateRepo")).Return(0, errors.New("Error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_missing_certificate_name",
			profileID:     "1",
			certificateID: "1",
			input: specs.UpdateCertificateRequest{
				Certificate: specs.Certificate{
					Name:             "",
					OrganizationName: "Organization",
					Description:      "Description",
					IssuedDate:       "2023-01-01",
					FromDate:         "2022-01-01",
					ToDate:           "2024-01-01",
				},
			},
			setup: func(certificateMock *mocks.CertificateStorer) {
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo")).Return(0, errors.New("Missing certificate name")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_invalid_profileid_or_certificateid",
			profileID:     "invalid",
			certificateID: "1",
			input: specs.UpdateCertificateRequest{
				Certificate: specs.Certificate{
					Name:             "Valid Name",
					OrganizationName: "Valid Organization",
					Description:      "Valid Description",
					IssuedDate:       "2023-01-01",
					FromDate:         "2022-01-01",
					ToDate:           "2024-01-01",
				},
			},
			setup:           func(certificateMock *mocks.CertificateStorer) {},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateRepo)

			_, err := certService.UpdateCertificate(context.TODO(), test.profileID, test.certificateID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestListCertificates(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	var repodeps = service.RepoDeps{
		CertificateDeps: mockCertificateRepo,
	}
	certificateService := service.NewServices(repodeps)

	mockProfileId := strconv.Itoa(profileID)
	mockResponseCertificate := []specs.CertificateResponse{
		{
			ProfileID:        123,
			Name:             "Certificate Name",
			OrganizationName: "Organization Name",
			Description:      "Certificate Description",
			IssuedDate:       "2024-01-01",
			FromDate:         "2024-01-01",
			ToDate:           "2024-06-01",
		},
	}

	tests := []struct {
		Name            string
		ProfileID       string
		MockSetup       func(*mocks.CertificateStorer, int)
		isErrorExpected bool
		wantResponse    []specs.CertificateResponse
	}{
		{
			Name:      "success_list_certificates",
			ProfileID: mockProfileId,
			MockSetup: func(mockCertificateStorer *mocks.CertificateStorer, profileID int) {
				mockCertificateStorer.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return(mockResponseCertificate, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseCertificate,
		},
		{
			Name:      "fail_get_certificates",
			ProfileID: "123",
			MockSetup: func(certMock *mocks.CertificateStorer, profileID int) {
				certMock.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{}, errors.New("error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.CertificateResponse{},
		},
		{
			Name:      "sucess_with_empty_resultset",
			ProfileID: "123",
			MockSetup: func(certMock *mocks.CertificateStorer, profileID int) {
				certMock.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{}, nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    []specs.CertificateResponse{},
		},
		{
			Name:      "invalid_profile_id",
			ProfileID: "invalid",
			MockSetup: func(certMock *mocks.CertificateStorer, profileID int) {
				certMock.On("ListCertificates", mock.Anything, mock.Anything, mock.Anything).Return([]specs.CertificateResponse{}, errors.New("invalid profile ID")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.CertificateResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			profileIDInt, _ := strconv.Atoi(tt.ProfileID)
			tt.MockSetup(mockCertificateRepo, profileIDInt)
			gotResponse, err := certificateService.ListCertificates(context.Background(), profileIDInt, specs.ListCertificateFilter{})
			assert.Equal(t, tt.wantResponse, gotResponse)
			if (err != nil) != tt.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", tt.Name, tt.isErrorExpected, err != nil)
			}
		})
	}
}
