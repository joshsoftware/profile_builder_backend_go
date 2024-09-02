package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
)

var (
	TestEmail = "test@example.com"
)

func TestUserLoginHandler(t *testing.T) {
	mockUserLoginService := new(mocks.Service)
	handlerFunc := handler.Login(context.Background(), mockUserLoginService)

	tests := []struct {
		Name               string
		Email              string
		AccessToken        string
		MockSendRequest    func(context.Context, string, string, string, io.Reader, map[string]string) ([]byte, error)
		MockSetup          func(*mocks.Service, string)
		RequestBody        specs.UserLoginRequest
		ExpectedStatusCode int
		ExpectedResponse   string
	}{
		{
			Name:        "success_of_login",
			Email:       TestEmail,
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				userInfo := specs.UserInfoFilter{ID: 0, Email: TestEmail}
				return json.Marshal(userInfo)
			},
			MockSetup: func(mockUserLoginService *mocks.Service, email string) {
				mockUserLoginService.On("GenerateLoginToken", context.Background(), specs.UserInfoFilter{Email: TestEmail}).Return(specs.LoginResponse{
					Token:     "valid_token",
					ProfileID: 1,
					Role:      "user",
				}, nil).Once()
			},
			RequestBody: specs.UserLoginRequest{
				AccessToken: "valid_access_token",
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"data":{"message":"Login successfully","profile_id":1,"role":"user","token":"valid_token","status_code":200}}`,
		},

		{
			Name:               "decoding_request_error",
			RequestBody:        specs.UserLoginRequest{},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error_code":400,"error_message":"empty access token"}`,
		},
		{
			Name:        "send_request_error",
			Email:       TestEmail,
			AccessToken: "invalid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				return nil, errors.New("failed to perform HTTP request")
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "invalid_access_token"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error_code":400,"error_message":"failed to perform HTTP request"}`,
		},
		{
			Name:        "decode_response_error",
			Email:       TestEmail,
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				return []byte("invalid json"), nil
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "valid_access_token"},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   `{"error_code":500,"error_message":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			Name:        "invalid_email_error",
			Email:       "",
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				userInfo := specs.UserInfoFilter{Email: ""}
				return json.Marshal(userInfo)
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "valid_access_token"},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"error_code":404,"error_message":"email not found"}`,
		},
		{
			Name:        "Fail_for_generate_login_token",
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				userInfo := specs.UserInfoFilter{ID: 0, Email: TestEmail}
				return json.Marshal(userInfo)
			},
			MockSetup: func(mockUserLoginService *mocks.Service, email string) {
				mockUserLoginService.On("GenerateLoginToken", context.Background(), specs.UserInfoFilter{Email: TestEmail}).Return(specs.LoginResponse{}, errs.ErrNoRecordFound).Once()
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "valid_access_token"},
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   `{"error_code":401,"error_message":"unauthorized access"}`,
		},
		{
			Name:        "service layer error",
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				userInfo := specs.UserInfoFilter{ID: 0, Email: TestEmail}
				return json.Marshal(userInfo)
			},
			MockSetup: func(mockUserLoginService *mocks.Service, email string) {
				mockUserLoginService.On("GenerateLoginToken", context.Background(), specs.UserInfoFilter{Email: TestEmail}).Return(specs.LoginResponse{}, errors.New("internal server error")).Once()
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "valid_access_token"},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   `{"error_code":500,"error_message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			var patchSendRequest *mpatch.Patch
			var err error

			if tt.MockSendRequest != nil {
				patchSendRequest, err = mpatch.PatchMethod(helpers.SendRequest, tt.MockSendRequest)
				if err != nil {
					t.Fatalf("Failed to patch SendRequest method: %v", err)
				}
				defer patchSendRequest.Unpatch()
			}

			if tt.MockSetup != nil {
				tt.MockSetup(mockUserLoginService, tt.Email)
			}

			reqBody, _ := json.Marshal(tt.RequestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			resp := httptest.NewRecorder()
			handlerFunc(resp, req)
			assert.Equal(t, tt.ExpectedStatusCode, resp.Code)

			if resp.Body.String() != tt.ExpectedResponse {
				t.Errorf("Expected response body %s but got %s", tt.ExpectedResponse, resp.Body.String())
			}

			body, _ := io.ReadAll(resp.Body)
			t.Log("response body is : ", string(body))
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	mockProfileService := new(mocks.Service)
	handlerFunc := handler.Logout(context.Background(), mockProfileService)

	tests := []struct {
		Name               string
		AuthHeader         string
		MockRemoveToken    func(token string) error
		ExpectedStatusCode int
		ExpectedResponse   string
	}{
		{
			Name:       "success_of_logout",
			AuthHeader: "Bearer valid_token",
			MockRemoveToken: func(token string) error {
				if token == "valid_token" {
					return nil
				}
				return errs.ErrTokenNotFound
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"data":{"message":"Logout successfully"}}`,
		},
		{
			Name:               "Fail_for_empty_token",
			AuthHeader:         "",
			MockRemoveToken:    func(token string) error { return nil },
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error_code":400,"error_message":"empty token"}`,
		},
		{
			Name:       "Fail_for_invalid_token",
			AuthHeader: "Bearer invalid_token",
			MockRemoveToken: func(token string) error {
				if token == "valid_token" {
					return nil
				}
				return errs.ErrTokenNotFound
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error_code":400,"error_message":"token not found in whitelist"}`,
		},
		{
			Name:       "Fail_for_remove_token_error",
			AuthHeader: "Bearer valid_token",
			MockRemoveToken: func(token string) error {
				return errors.New("failed to remove token")
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error_code":400,"error_message":"failed to remove token"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockProfileService.On("RemoveToken", strings.TrimPrefix(test.AuthHeader, "Bearer ")).Return(test.MockRemoveToken(strings.TrimPrefix(test.AuthHeader, "Bearer "))).Once()

			req := httptest.NewRequest("POST", "/logout", nil)
			if test.AuthHeader != "" {
				req.Header.Set("Authorization", test.AuthHeader)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlerFunc)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.ExpectedStatusCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != test.ExpectedResponse {
				t.Errorf("Expected response body %s, got %s", test.ExpectedResponse, rr.Body.String())
			}
		})
	}
}
