package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
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
				userInfo := specs.UserInfoFilter{Email: TestEmail}
				return json.Marshal(userInfo)
			},
			MockSetup: func(mockUserLoginService *mocks.Service, email string) {
				mockUserLoginService.On("GenerateLoginToken", context.Background(), email).Return("valid_token", nil).Once()
			},
			RequestBody: specs.UserLoginRequest{
				AccessToken: "valid_access_token",
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"data":{"message":"Login successful","token":"valid_token","status_code":200}}`,
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
			Name:        "generate_token_error",
			Email:       TestEmail,
			AccessToken: "valid_access_token",
			MockSendRequest: func(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
				userInfo := specs.UserInfoFilter{Email: TestEmail}
				return json.Marshal(userInfo)
			},
			MockSetup: func(mockUserLoginService *mocks.Service, email string) {
				mockUserLoginService.On("GenerateLoginToken", context.Background(), email).Return("", errors.New("Unable to generate token")).Once()
			},
			RequestBody:        specs.UserLoginRequest{AccessToken: "valid_access_token"},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   `{"error_code":500,"error_message":"Unable to generate token"}`,
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
			respBody := resp.Body.String()
			t.Log("response is : ", respBody)

			if !reflect.DeepEqual(resp.Body.String(), tt.ExpectedResponse) {
				t.Errorf("\n HandlerFunc() = %s, \n\n\n want %s, diff:%+v", resp.Body.String(), tt.ExpectedResponse, cmp.Diff(resp.Body.String(), tt.ExpectedResponse))
			}
		})
	}
}
