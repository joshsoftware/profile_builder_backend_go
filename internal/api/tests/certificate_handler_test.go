package test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
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

func TestGetCertificatesHandler(t *testing.T) {
	certficateSvc := new(mocks.Service)
	getCertificateHandler := handler.GetCertificatesHandler(context.Background(), certficateSvc)

	tests := []struct {
		name               string
		queryParams        string
		mockSvcSetup       func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success for fetching certificate details",
			queryParams: profileID,
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCertificates", mock.Anything, profileID).Return([]dto.CertificateResponse{
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
			name:        "Success for fetching multiple certificates",
			queryParams: profileID,
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCertificates", mock.Anything, profileID).Return([]dto.CertificateResponse{
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
			name:        "Fail as error in GetCertificates",
			queryParams: profileID,
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCertificates", mock.Anything, profileID).Return([]dto.CertificateResponse{}, errors.New("some error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:        "Success for fetching no certificates",
			queryParams: profileID,
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCertificates", mock.Anything, profileID).Return([]dto.CertificateResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Fail as error in GetParams",
			queryParams: profileID0,
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCertificates", mock.Anything, profileID0).Return(nil, errors.New("invalid profile id")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSvcSetup(certficateSvc)

			req, err := http.NewRequest("GET", "/profiles/"+tt.queryParams+"/certificates", nil)
			if err != nil {
				t.Fatal(err)
				return
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.queryParams})

			resp := httptest.NewRecorder()
			handler := http.HandlerFunc(getCertificateHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("Expected %d but got %d", tt.expectedStatusCode, resp.Result().StatusCode)
			}
		})
	}
}
