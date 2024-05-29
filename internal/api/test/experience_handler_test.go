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

func TestCreateExperienceHandler(t *testing.T) {
	ctx := context.Background()
	profileSvc := new(mocks.Service)
	createExperienceHandler := handler.CreateExperienceHandler(ctx, profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for experience detail",
			input: `{
				"profile_id": 1,
				"experiences":[{
					"designation": "Associate Data Scientist",
					"company_name": "Josh Software Pvt.Ltd.",
					"from_date": "Jan-2023",
					"to_date": "July-2024"
					}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateExperience", context.Background(), mock.AnythingOfType("dto.CreateExperienceRequest")).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect JSON",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing designation field",
			input: `{
                "profile_id": 1,
                "experiences": [{
                    "designation": "",
                    "company_name": "ABC Corp",
                    "from_date": "2023-01-01",
                    "to_date": "2024-01-01"
                }]
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
                "profile_id": 0,
                "experiences": [{
                    "designation": "Software Engineer",
                    "company_name": "ABC Corp",
                    "from_date": "2023-01-01",
                    "to_date": "2024-01-01"
                }]
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// mockSvc := new(mocks.Service)
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/experiences", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
