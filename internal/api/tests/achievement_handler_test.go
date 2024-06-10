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

var (
	profileID  = 1
	profileID0 = 0
)

func TestListAchievementsHandler(t *testing.T) {
	achSvc := new(mocks.Service)
	getAchievementHandler := handler.ListAchievementsHandler(context.Background(), achSvc)

	tests := []struct {
		name               string
		queryParams        int
		MockSvc            func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "success_achievement",
			queryParams: profileID,
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID).Return([]dto.AchievementResponse{
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
			queryParams: profileID,
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID).Return([]dto.AchievementResponse{
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
			queryParams: profileID,
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID).Return([]dto.AchievementResponse{}, errors.New("some error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
		{
			name:        "sucess_with_empty_resultset",
			queryParams: profileID,
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID).Return([]dto.AchievementResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{

			name:        "fail_to_fetch_achievements_with_invalid_profile_id",
			queryParams: profileID0,
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID0).Return(nil, errors.New("invalid profile id")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockSvc(achSvc)
			req, err := http.NewRequest("GET", "/profiles/"+strconv.Itoa(tt.queryParams)+"/achievements", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": strconv.Itoa(tt.queryParams)})

			resp := httptest.NewRecorder()

			handler := http.HandlerFunc(getAchievementHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, resp.Code)
			}
		})
	}

}
