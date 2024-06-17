package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				    }]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("specs.CreateAchievementRequest"), 1).Return(1, nil).Once()
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
				"achievements":[{
				    "name": "",
					"description": "Description of Award"
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing description field",
			input: `{
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
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestUpdateAchievementHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	updateAchievementHandler := handler.UpdateAchievementHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_updating_achievement_detail",
			input: `{
				"achievement": {
					"name": "Updated Star Performer",
					"description": "Updated Description"
				}
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateAchievement", context.Background(), "1", "1", mock.AnythingOfType("specs.UpdateAchievementRequest")).Return(1, nil).Once()
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
			name: "Fail_fo_missing_name_field",
			input: `{
				"achievement": {
					"name": "",
					"description": "Updated Description"
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_description_field",
			input: `{
				"achievement": {
					"name": "Updated Star Performer",
					"description": ""
				}
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("PUT", "/profiles/1/achievements/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
