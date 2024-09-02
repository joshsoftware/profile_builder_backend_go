package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
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
		expectedResponse   string
	}{
		{
			name: "Success for valid achievement details",
			input: `{
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				  }]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("specs.CreateAchievementRequest"), 1, 1).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Achievement(s) added successfully","profile_id":1}}`,
		},
		{
			name: "Success_for_multiple_achievements",
			input: `{
				"achievements":[
					{
						"name": "Star Performer",
						"description": "First Achievement"
					},
					{
						"name": "Best Developer",
						"description": "Second Achievement"
					}
				]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("specs.CreateAchievementRequest"), 1, 1).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Achievement(s) added successfully","profile_id":1}}`,
		},
		{
			name: "Success_without_description",
			input: `{
				"achievements":[{
					"name": "Star Performer",
					"description": ""
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("specs.CreateAchievementRequest"), 1, 1).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"data":{"message":"Achievement(s) added successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
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
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : name "}`,
		},
		{
			name: "Fail_when_CreateAchievement_fails",
			input: `{
				"achievements":[{
					"name": "Star Performer",
					"description": "Description of Award"
				}]
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("specs.CreateAchievementRequest"), 1, 1).Return(0, errors.New("service failure")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"service failure"}`,
		},
		{
			name: "Fail for invalid profile ID",
			input: `{
				"achievements":[{
					"name": "Star Performer",
					"description": "Description of Award"
				}]
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid request data"}`,
		},
		{
			name: "Fail for missing UserID in context",
			input: `{
				"achievements":[{
					"name": "Star Performer",
					"description": "Description of Award"
				}]
			}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid user id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)
			defer profileSvc.AssertExpectations(t)
			req := httptest.NewRequest("POST", "/profiles/achievements", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1"})
			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			if test.name == "Fail for invalid profile ID" {
				req = mux.SetURLVars(req, map[string]string{"profile_id": "invalid"})
			}

			if test.name == "Fail for missing UserID in context" {
				ctx := context.WithValue(req.Context(), constants.UserIDKey, 1)
				req = req.WithContext(ctx)
			}

			handler := http.HandlerFunc(createAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.expectedResponse {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
			}
		})
	}
}

func TestUpdateAchievementHandler(t *testing.T) {
	achSvc := new(mocks.Service)
	updateAchievementHandler := handler.UpdateAchievementHandler(context.Background(), achSvc)
	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success_for_achievement_update",
			input: `{
						"achievement":{
    						"name": "Star Performer of JOSH",
							"description": "Description of Award"
    					}
					}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateAchievement", mock.Anything, 1, 1, 1, mock.AnythingOfType("specs.UpdateAchievementRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Achievement updated successfully","profile_id":1}}`,
		},
		{
			name: "Success_without_description",
			input: `{
						"achievement":{
							"name": "Star Performer of JOSH",
							"description": ""
						}
					}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateAchievement", mock.Anything, 1, 1, 1, mock.AnythingOfType("specs.UpdateAchievementRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Achievement updated successfully","profile_id":1}}`,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid request body"}`,
		},
		{
			name: "Fail_for_missing_name_field",
			input: `{
						"achievement":{
							"name": "",
							"description": "Description of Award"
						}
					}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"parameter missing : name "}`,
		},
		{
			name: "Fail_for_service_error",
			input: `{
						"achievement":{
							"name": "Star Performer of JOSH",
							"description": "Description of Award"
						}
					}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateAchievement", mock.Anything, 1, 1, 1, mock.AnythingOfType("specs.UpdateAchievementRequest")).Return(0, errors.New("service error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"service error"}`,
		},
		{
			name: "Fail_for_invalid_profile_id",
			input: `{
						"achievement":{
							"name": "Star Performer of JOSH",
							"description": "Description of Award"
						}
					}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid request data"}`,
		},
		{
			name: "Fail_for_missing_user_id_in_context",
			input: `{
						"achievement":{
							"name": "Star Performer of JOSH",
							"description": "Description of Award"
						}
					}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid user id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(achSvc)

			req := httptest.NewRequest("PUT", "/profiles/achievements/1", bytes.NewBuffer([]byte(test.input)))
			req = mux.SetURLVars(req, map[string]string{"profile_id": "1", "id": "1"})

			ctx := context.WithValue(req.Context(), constants.UserIDKey, 1.0)
			req = req.WithContext(ctx)

			if test.name == "Fail_for_invalid_profile_id" {
				req = mux.SetURLVars(req, map[string]string{"profile_id": "invalid", "id": "1"})
			}

			if test.name == "Fail_for_missing_user_id_in_context" {
				ctx := context.WithValue(req.Context(), constants.UserIDKey, 1)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.expectedResponse {
				t.Errorf("Expected response body %s but got %s", test.expectedResponse, rr.Body.String())
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
		expectedResponse   string
	}{
		{
			name:        "success_achievement",
			pathParams:  profileID,
			queryParams: "achievement_ids=1,2&names=Client%20Appreciation",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{
						AchievementIDs: []int{1},
						Names:          []string{"Client Appreciation"},
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
			expectedResponse:   `{"data":{"achievements":[{"id":0,"profile_id":1,"name":"Client Appreciation","description":"Description of Appreciation"}]}}`,
		},
		{
			name:        "success_achievements",
			pathParams:  profileID,
			queryParams: "achievement_ids=1,2&achievement_names=Client%20Appreciation,Another%20Achievement",
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
			expectedResponse:   `{"data":{"achievements":[{"id":0,"profile_id":1,"name":"Client Appreciation","description":"Description of Appreciation"},{"id":0,"profile_id":1,"name":"Another Achievement","description":"Description of Another Achievement"}]}}`,
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
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"achievements":[]}}`,
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
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
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
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:        "failed_because_service_layer_caused_error",
			pathParams:  profileID,
			queryParams: "achievement_ids=1",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{
						AchievementIDs: []int{1, 2},
						Names:          []string{"Client Appreciation", "Another Achievement"},
					}, nil
				})
			},
			MockSvc: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAchievements", mock.Anything, profileID, mock.Anything).Return([]specs.AchievementResponse{}, errors.New("serive layer error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"failed to fetch data"}`,
		},
		{
			name:               "invalid_profile_id",
			mockDecodeRequest:  func() {},
			MockSvc:            func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   `{"error_code":502,"error_message":"invalid profile id"}`,
		},
		{
			name:        "failed_to_decode_request",
			pathParams:  profileID,
			queryParams: "achievement_ids=a",
			mockDecodeRequest: func() {
				mpatch.PatchMethod(helpers.DecodeAchievementRequest, func(r *http.Request) (specs.ListAchievementFilter, error) {
					return specs.ListAchievementFilter{}, errors.New("failed to decode request")
				})
			},
			MockSvc:            func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"unable to decode request"}`, // Adjust error message & code here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockSvc(achSvc)
			req := httptest.NewRequest("GET", "/profiles/"+strconv.Itoa(tt.pathParams)+"/achievements"+tt.queryParams, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": strconv.Itoa(tt.pathParams)})

			if tt.name == "invalid_profile_id" {
				req = mux.SetURLVars(req, map[string]string{})
			}

			if tt.name == "failed_to_decode_request" {
				tt.mockDecodeRequest()
			}
			resp := httptest.NewRecorder()
			handler := http.HandlerFunc(getAchievementHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, resp.Code)
			}

			if resp.Body.String() != tt.expectedResponse {
				t.Errorf("Expected response body %s but got %s, \n\n ######diff :%+v", tt.expectedResponse, resp.Body.String(), cmp.Diff(tt.expectedResponse, resp.Body.String()))
			}

			body := resp.Body.String()
			fmt.Println(body)
		})
	}

}

func TestDeleteAchievementHandler(t *testing.T) {
	achSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		achievementID      string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:          "Success_for_achievement_delete",
			profileID:     "1",
			achievementID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteAchievement", mock.Anything, 1, 1).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Achievement deleted successfully",
		},
		{
			name:          "No_data_found_for_deletion",
			profileID:     "1",
			achievementID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteAchievement", mock.Anything, 1, 2).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Resource not found for the given request ID"}}`,
		},
		{
			name:          "Error_while_deleting_achievement",
			profileID:     "1",
			achievementID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteAchievement", mock.Anything, 1, 3).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			achievementID:      "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
		{
			name:          "Fail_for_internal_error",
			profileID:     "1",
			achievementID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteAchievement", mock.Anything, 1, 1).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(achSvc)
			reqPath := "/profiles/" + tt.profileID + "/achievements/" + tt.achievementID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID, "id": tt.achievementID})
			rr := httptest.NewRecorder()

			handler := handler.DeleteAchievementHandler(context.Background(), achSvc)
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
