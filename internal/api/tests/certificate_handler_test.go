package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
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
