package test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
)

func TestCreateCertificateHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	createCertificateHandler := handler.CreateCertificateHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for certificate Detail",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCertificate", mock.Anything, mock.AnythingOfType("dto.CreateCertificateRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing description field",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/certificates", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createCertificateHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestListCertificatesHandler(t *testing.T) {
	certficateSvc := new(mocks.Service)
	getCertificateHandler := handler.ListCertificatesHandler(context.Background(), certficateSvc)

	tests := []struct {
		name               string
		pathParams         int
		queryParams        string
		mockDecodeRequest  func()
		mockSvcSetup       func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success_for_fetching_single_certificate",
			pathParams:  profileID,
			queryParams: "certificate_ids=1,2&names=Golang,ROR",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (dto.ListCertificateFilter, error) {
					return dto.ListCertificateFilter{
						CertificateIDs: []int{1, 2},
						Names:          []string{"Golang", "ROR"},
					}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]dto.CertificateResponse{
					{
						ProfileID:        1,
						Name:             "Golang Master Class",
						OrganizationName: "Udemy",
						Description:      "A Bootcamp for Mastering Golang Concepts",
						IssuedDate:       "Dec-2023",
						FromDate:         "June-2023",
						ToDate:           "Dec-2023",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "success_for_fetching_multiple_certificates",
			pathParams:  profileID,
			queryParams: "certificate_ids=1,2&names=Golang,ROR",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (dto.ListCertificateFilter, error) {
					return dto.ListCertificateFilter{
						CertificateIDs: []int{1, 2},
						Names:          []string{"Golang", "ROR"},
					}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]dto.CertificateResponse{
					{
						ProfileID:   1,
						Name:        "Certificate 1",
						Description: "Description of Certificate 1",
					},
					{
						ProfileID:   1,
						Name:        "Certificate 2",
						Description: "Description of Certificate 2",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},

		{
			name:        "fail_to_fetch_certificates",
			pathParams:  profileID,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (dto.ListCertificateFilter, error) {
					return dto.ListCertificateFilter{}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]dto.CertificateResponse{}, errors.New("some error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:        "sucess_with_empty_resultset",
			pathParams:  profileID,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (dto.ListCertificateFilter, error) {
					return dto.ListCertificateFilter{}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]dto.CertificateResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "fail_to_fetch_certificates_with_invalid_profile_id",
			pathParams:  profileID0,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (dto.ListCertificateFilter, error) {
					return dto.ListCertificateFilter{}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID0, mock.Anything).Return(nil, errors.New("invalid profile id")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSvcSetup(certficateSvc)

			req, err := http.NewRequest("GET", "/profiles/"+strconv.Itoa(tt.pathParams)+"/certificates", nil)
			if err != nil {
				t.Fatal(err)
				return
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": strconv.Itoa(tt.pathParams)})
			resp := httptest.NewRecorder()

			handler := http.HandlerFunc(getCertificateHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("Expected %d but got %d", tt.expectedStatusCode, resp.Result().StatusCode)
			}
		})
	}
}
