package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
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
			name: "Success_for_certificate_detail",
			input: `{
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
				mockSvc.On("CreateCertificate", mock.Anything, mock.AnythingOfType("specs.CreateCertificateRequest"), 1).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_name_field",
			input: `{
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
			name: "Fail_for_missing_description_field",
			input: `{
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

			req, err := http.NewRequest("POST", "/profiles/1/certificates", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})
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
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (specs.ListCertificateFilter, error) {
					return specs.ListCertificateFilter{
						CertificateIDs: []int{1, 2},
						Names:          []string{"Golang", "ROR"},
					}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{
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
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (specs.ListCertificateFilter, error) {
					return specs.ListCertificateFilter{
						CertificateIDs: []int{1, 2},
						Names:          []string{"Golang", "ROR"},
					}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{
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
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (specs.ListCertificateFilter, error) {
					return specs.ListCertificateFilter{}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{}, errors.New("some error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:        "sucess_with_empty_resultset",
			pathParams:  profileID,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (specs.ListCertificateFilter, error) {
					return specs.ListCertificateFilter{}, nil
				})
			},
			mockSvcSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCertificates", mock.Anything, profileID, mock.Anything).Return([]specs.CertificateResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:        "fail_to_fetch_certificates_with_invalid_profile_id",
			pathParams:  profileID0,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeCertificateRequest, func(r *http.Request) (specs.ListCertificateFilter, error) {
					return specs.ListCertificateFilter{}, nil
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

func TestUpdateCertificateHandler(t *testing.T) {
	certificateSvc := new(mocks.Service)
	updateCertificateHandler := handler.UpdateCertificateHandler(context.Background(), certificateSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_updating_certificate_detail",
			input: `{
				"certificate": {
					"name": "Updated Certificate",
					"organization_name": "Updated Organization",
					"description": "Updated Description",
					"issued_date": "2024-05-30",
					"from_date": "2023-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCertificate", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateCertificateRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_name_field",
			input: `{
				"certificate": {
					"name": "",
					"organization_name": "Updated Organization",
					"description": "Updated Description",
					"issued_date": "2024-05-30",
					"from_date": "2023-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_organization_name_field",
			input: `{
				"certificate": {
					"name": "Updated Certificate",
					"organization_name": "",
					"description": "Updated Description",
					"issued_date": "2024-05-30",
					"from_date": "2023-01-01",
					"to_date": "2023-12-31"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(certificateSvc)

			req, err := http.NewRequest("PUT", "/profiles/1/certificates/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateCertificateHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestDeleteCertificateHandler(t *testing.T) {
	certificateSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		certificateID      string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:          "Success_for_certificate_delete",
			profileID:     "1",
			certificateID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCertificate", mock.Anything, mock.AnythingOfType("specs.DeleteCertificateRequest")).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Certificate deleted successfully",
		},
		{
			name:          "No_data_found_for_deletion",
			profileID:     "1",
			certificateID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCertificate", mock.Anything, mock.AnythingOfType("specs.DeleteCertificateRequest")).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "No data found for deletion",
		},
		{
			name:          "Error_while_deleting_certificate",
			profileID:     "1",
			certificateID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCertificate", mock.Anything, mock.AnythingOfType("specs.DeleteCertificateRequest")).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			certificateID:      "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
		{
			name:          "Fail_for_internal_error",
			profileID:     "1",
			certificateID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCertificate", mock.Anything, mock.AnythingOfType("specs.DeleteCertificateRequest")).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(certificateSvc)
			reqPath := "/profiles/" + tt.profileID + "/certificate/" + tt.certificateID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID, "id": tt.certificateID})
			rr := httptest.NewRecorder()

			handler := handler.DeleteCertificatesHandler(context.Background(), certificateSvc)
			handler(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if !strings.Contains(string(body), tt.expectedResponse) {
				t.Errorf("expected response to contain %q, got %q", tt.expectedResponse, body)
			}
		})
	}

}
