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
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
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
			name: "Success_for_achievement_detail",
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
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_name_field",
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
			name: "Fail_for_missing_description_field",
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

			ctx := context.WithValue(req.Context(), userIDKey, 1)
			req = req.WithContext(ctx)

			handler := http.HandlerFunc(createAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

// Define a custom type for context key
type contextKey string

// Define constants for context keys
const (
	userIDKey        contextKey = "user_id"
	profileIDKey     contextKey = "profile_id"
	achievementIDKey contextKey = "achievement_id"
)

func TestUpdateAchievementHandler(t *testing.T) {
	achSvc := new(mocks.Service)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_achievement_update",
			input: `{
				"name": "Updated Star Performer",
				"description": "Updated description of Award"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateAchievement", mock.Anything, 1, 1, 1, mock.AnythingOfType("specs.UpdateAchievementRequest")).Return(1, nil).Once()
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
				"name": "",
				"description": "Updated description of Award"
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_description_field",
			input: `{
				"name": "Updated Star Performer",
				"description": ""
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(achSvc)

			req, err := http.NewRequest("PUT", "/profiles/achievements/1", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "achievement_id": "1"})

			// Set the user_id in the context
			ctx := context.WithValue(req.Context(), userIDKey, 1)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.UpdateAchievementHandler(ctx, achSvc)(w, r)
			})
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

var (
	profileID  = 1
	profileID0 = 0
)

func TestListAchievementsHandler(t *testing.T) {
	achSvc := new(mocks.Service)
	getAchievementHandler := handler.ListAchievementsHandler(context.Background(), achSvc)

	tests := []struct {
		name               string
		pathParams         int
		queryParams        string
		mockDecodeRequest  func()
		MockSvc            func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "success_achievement",
			pathParams:  profileID,
			queryParams: "achievement_ids=1,2&names=Client Appreciation,Another Achievement",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{
						AchievementIDs: []int{1, 2},
						Names:          []string{"Client Appreciation", "Another Achievement"},
					}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{
					{
						ProfileID:   1,
						Name:        "Client Appreciation",
						Description: "Description of Appreciation",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "success_achievements",
			pathParams:  profileID,
			queryParams: "achievement_ids=1,2&achievement_names=Client Appreciation,Another Achievement",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{
						AchievementIDs: []int{1, 2},
						Names:          []string{"Client Appreciation", "Another Achievement"},
					}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{
					{
						ProfileID:   1,
						Name:        "Client Appreciation",
						Description: "Description of Appreciation",
					},
					{
						ProfileID:   1,
						Name:        "Another Achievement",
						Description: "Description of Another Achievement",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "fail_to_fetch_achievements",
			pathParams:  profileID,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("some error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:        "sucess_with_empty_resultset",
			pathParams:  profileID,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:        "fail_to_fetch_achievements_with_invalid_profile_id",
			pathParams:  profileID0,
			queryParams: "",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID0, mock.Anything).Return(nil, errors.New("invalid profile id")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockSvc(achSvc)
			req, err := http.NewRequest("GET", "/profiles/"+strconv.Itoa(tt.pathParams)+"/achievements", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": strconv.Itoa(tt.pathParams)})
			resp := httptest.NewRecorder()

			handler := http.HandlerFunc(getAchievementHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, resp.Code)
			}
		})
	}

}
