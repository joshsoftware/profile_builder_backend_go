package post_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/profile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	userSrv := mocks.NewService(t)
	loginHandler := post.Login(context.Background(), userSrv)

	tests := []struct {
		name               string
		requestBody        string
		setup              func(srvMock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Successful login",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				os.Setenv("GOOGLE_USER_INFO_URL", "https://valid.url")
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte(`{"email":"user@example.com"}`))),
						}, nil
					}),
				}
				srvMock.On("GenerateLoginToken", mock.Anything, "user@example.com").Return("valid-login-token", nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "DecodeUserLoginRequest error",
			requestBody:        "{invalid-json}",
			setup:              func(srvMock *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "NewRequest Error",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				os.Setenv("GOOGLE_USER_INFO_URL", "/someinvalidurl")
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
						}, fmt.Errorf("some error")
					}),
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "HTTP Client Do Error",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return nil, errors.New("HTTP request error")
					}),
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "readAll error",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(errReader(0)),
						}, nil
					}),
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "json decode error",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte("invalid-json"))),
						}, nil
					}),
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "GenerateLoginToken error",
			requestBody: `{"access_token":"valid-token"}`,
			setup: func(srvMock *mocks.Service) {
				os.Setenv("GOOGLE_USER_INFO_URL", "https://valid.url")
				http.DefaultClient = &http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte(`{"email":"user@example.com"}`))),
						}, nil
					}),
				}
				userSrv.On("GenerateLoginToken", mock.Anything, "user@example.com").Return("", errors.New("token generation error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setup(userSrv)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(tt.requestBody)))
			w := httptest.NewRecorder()

			loginHandler(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
