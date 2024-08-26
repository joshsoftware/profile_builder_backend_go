package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCertificate(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
		CertificateDeps: mockCertificateRepo,
	}
	certificateService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		input           specs.CreateCertificateRequest
		profileId       int
		userID          int
		setup           func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer)
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
			profileId: 1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name: "failed_because_count_records_returns_an_error",
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
			profileId: 1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_createcertificate_6greturns_an_error",
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
			profileId: 1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
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
			profileId: 1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(errors.New("Missing certificate name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_empty_payload",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{{}},
			},
			profileId: 1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(errors.New("Empty payload")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_invalid_profile_id",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{
					{
						Name:             "Certificate D",
						OrganizationName: "Organization D",
						Description:      "Description D",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			profileId: -1,
			userID:    1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(errors.New("Invalid profile ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name: "Failed_because_of_invalid_user_id",
			input: specs.CreateCertificateRequest{
				Certificates: []specs.Certificate{
					{
						Name:             "Certificate E",
						OrganizationName: "Organization E",
						Description:      "Description E",
						IssuedDate:       "2024-01-01",
						FromDate:         "2024-01-01",
						ToDate:           "2024-06-01",
					},
				},
			},
			profileId: 1,
			userID:    -1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				profileMock.On("CountRecords", mock.Anything, 1, mock.Anything, mock.Anything).Return(1, nil).Once()
				certificateMock.On("CreateCertificate", mock.Anything, mock.AnythingOfType("[]repository.CertificateRepo"), mock.Anything).Return(errors.New("Invalid user ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateRepo, mockProfileRepo)
			_, err := certificateService.CreateCertificate(context.TODO(), test.input, 1, 1)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}

			mockProfileRepo.AssertExpectations(t)
			mockCertificateRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateCertificate(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
		CertificateDeps: mockCertificateRepo,
	}
	certService := service.NewServices(repodeps)

	tests := []struct {
		name            string
		profileID       int
		certificateID   int
		userID          int
		input           specs.UpdateCertificateRequest
		setup           func(*mocks.CertificateStorer, *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_updating_certificate_details",
			profileID:     1,
			certificateID: 1,
			userID:        1,
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
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo"), mock.Anything).Return(1, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed_because_updatecertificate_returns_an_error",
			profileID:     1,
			certificateID: 1,
			userID:        1,
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
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo"), mock.Anything).Return(0, errors.New("Error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_missing_certificate_name",
			profileID:     1,
			certificateID: 1,
			userID:        1,
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
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo"), mock.Anything).Return(0, errors.New("Missing certificate name")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_invalid_profileid_or_certificateid",
			profileID:     -1,
			certificateID: 1,
			userID:        1,
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
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("UpdateCertificate", mock.Anything, -1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo"), mock.Anything).Return(0, errors.New("Invalid profile ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_of_invalid_userid",
			profileID:     1,
			certificateID: 1,
			userID:        -1,
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
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("UpdateCertificate", mock.Anything, 1, 1, mock.AnythingOfType("repository.UpdateCertificateRepo"), mock.Anything).Return(0, errors.New("Invalid user ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateRepo, mockProfileRepo)

			_, err := certService.UpdateCertificate(context.TODO(), test.profileID, test.certificateID, test.userID, test.input)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}

			mockProfileRepo.AssertExpectations(t)
			mockCertificateRepo.AssertExpectations(t)
		})
	}
}

func TestListCertificates(t *testing.T) {
	mockCertificateRepo := new(mocks.CertificateStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repodeps = service.RepoDeps{
		ProfileDeps:     mockProfileRepo,
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
		MockSetup       func(*mocks.CertificateStorer, *mocks.ProfileStorer, int)
		isErrorExpected bool
		wantResponse    []specs.CertificateResponse
	}{
		{
			Name:      "success_list_certificates",
			ProfileID: mockProfileId,
			MockSetup: func(mockCertificateStorer *mocks.CertificateStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				mockCertificateStorer.On("ListCertificates", mock.Anything, profileID, mock.Anything, mock.Anything).Return(mockResponseCertificate, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    mockResponseCertificate,
		},
		{
			Name:      "fail_get_certificates",
			ProfileID: "123",
			MockSetup: func(certMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certMock.On("ListCertificates", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.CertificateResponse{}, errors.New("error")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.CertificateResponse{},
		},
		{
			Name:      "sucess_with_empty_resultset",
			ProfileID: "123",
			MockSetup: func(certMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certMock.On("ListCertificates", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.CertificateResponse{}, nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
			wantResponse:    []specs.CertificateResponse{},
		},
		{
			Name:      "invalid_profile_id",
			ProfileID: "invalid",
			MockSetup: func(certMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer, profileID int) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certMock.On("ListCertificates", mock.Anything, profileID, mock.Anything, mock.Anything).Return([]specs.CertificateResponse{}, errors.New("invalid profile ID")).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("handle transaction error")).Once()
			},
			isErrorExpected: true,
			wantResponse:    []specs.CertificateResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			profileIDInt, _ := strconv.Atoi(tt.ProfileID)
			tt.MockSetup(mockCertificateRepo, mockProfileRepo, profileIDInt)
			gotResponse, err := certificateService.ListCertificates(context.Background(), profileIDInt, specs.ListCertificateFilter{})
			assert.Equal(t, tt.wantResponse, gotResponse)
			if (err != nil) != tt.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", tt.Name, tt.isErrorExpected, err != nil)
			}
			mockProfileRepo.AssertExpectations(t)
			mockCertificateRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCertificateService(t *testing.T) {
	mockCertificateSvc := new(mocks.CertificateStorer)
	mockProfileRepo := new(mocks.ProfileStorer)
	var repoDeps = service.RepoDeps{
		CertificateDeps: mockCertificateSvc,
		ProfileDeps:     mockProfileRepo,
	}
	certificateSvc := service.NewServices(repoDeps)

	tests := []struct {
		name            string
		certificateID   int
		profileID       int
		setup           func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer)
		isErrorExpected bool
	}{
		{
			name:          "Success_for_delete_certificate",
			certificateID: 1,
			profileID:     1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("DeleteCertificate", mock.Anything, 1, 1, nil).Return(nil).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:          "Failed_because_delete_certificate_returns_an_error",
			certificateID: 2,
			profileID:     1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("DeleteCertificate", mock.Anything, 1, 2, nil).Return(errs.ErrNoData).Once()
				profileMock.On("HandleTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:          "Failed_because_DeleteCertificate_returns_an_error",
			certificateID: 3,
			profileID:     1,
			setup: func(certificateMock *mocks.CertificateStorer, profileMock *mocks.ProfileStorer) {
				profileMock.On("BeginTransaction", mock.Anything).Return(nil, nil).Once()
				certificateMock.On("DeleteCertificate", mock.Anything, 1, 3, nil).Return(errs.ErrFailedToDelete).Once()
				profileMock.On("HandleTransaction", mock.Anything, nil, mock.Anything).Return(nil).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockCertificateSvc, mockProfileRepo)
			err := certificateSvc.DeleteCertificate(context.Background(), test.profileID, test.certificateID)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test %s failed, expected error to be %v, but got err %v", test.name, test.isErrorExpected, err != nil)
			}
		})
	}
}
