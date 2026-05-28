package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/mock"
)

func TestInviteAdminHandler(t *testing.T) {
	mockService := new(mocks.Service)
	handlerFunc := handler.InviteAdmin(context.Background(), mockService)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		ctxSetup           func(r *http.Request) *http.Request
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success_for_admin_invite",
			input: `{
				"name": "Bruce Wayne",
				"email": "batman@gotham.com"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("InviteAdmin", mock.Anything, 1, specs.AdminInviteRequest{
					Name:  "Bruce Wayne",
					Email: "batman@gotham.com",
				}).Return(nil).Once()
			},
			ctxSetup: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), constants.UserIDKey, float64(1))
				return r.WithContext(ctx)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"message":"Admin invited successfully"}}`,
		},
		{
			name: "Fail_due_to_invalid_user_id",
			input: `{
				"name": "Bruce Wayne",
				"email": "batman@gotham.com"
			}`,
			setup: func(mockSvc *mocks.Service) {
			},
			ctxSetup: func(r *http.Request) *http.Request {
				// No UserIDKey in context
				return r
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error_code":400,"error_message":"invalid user id"}`,
		},
		{
			name: "Fail_due_to_duplicate_email",
			input: `{
				"name": "Bruce Wayne",
				"email": "batman@gotham.com"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("InviteAdmin", mock.Anything, 1, specs.AdminInviteRequest{
					Name:  "Bruce Wayne",
					Email: "batman@gotham.com",
				}).Return(errors.ErrDuplicateKey).Once()
			},
			ctxSetup: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), constants.UserIDKey, float64(1))
				return r.WithContext(ctx)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `{"error_code":409,"error_message":"record already exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockService)

			req := httptest.NewRequest(http.MethodPost, "/admin_invite", bytes.NewBuffer([]byte(tt.input)))
			req = tt.ctxSetup(req)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlerFunc)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected %d but got %d", tt.expectedStatusCode, rr.Result().StatusCode)
			}

			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("Expected response body %s but got %s", tt.expectedResponse, rr.Body.String())
			}
		})
	}
}
