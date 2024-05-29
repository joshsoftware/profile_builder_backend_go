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

func TestCreateAchievementHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	createAchievementHandler := handler.CreateAchievementHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for achievement detail",
			input: `{
				"profile_id": 1,
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				    }]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("dto.CreateAchievementRequest")).Return(1, nil).Once()
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
				"achievements":[{
				    "name": "",
					"description": "Description of Award"
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing description field",
			input: `{
				"profile_id": 1,
				"achievements":[{
				    "name": "Star Performer",
					"description": ""
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/achievements", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
