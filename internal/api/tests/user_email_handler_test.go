package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	service *mocks.Service
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.service = &mocks.Service{}
}

func (suite *HandlerTestSuite) TearDownSuite() {
	suite.service.AssertExpectations(suite.T())
}

func (s *HandlerTestSuite) TestSendUserInvitation() {
	t := s.T()

	type args struct {
		userID    string
		profileID string
		ctx       context.Context
	}

	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantResponse string
		prepare      func(args)
	}{
		// POSITIVE || Successfully send user invitation
		{
			name: "Successfully send user invitation",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"data":{"message":"Invitation sent successfully to employee"}}`,
			wantStatus:   http.StatusOK,
			prepare: func(a args) {
				s.service.On("SendUserInvitation", mock.Anything, TestUserID, TestProfileID).Return(nil).Once()
			},
		},
		// NEGATIVE || Fail due to invalid user id
		{
			name: "Fail_due_to_invalid_user_id",
			args: args{
				userID:    "",
				profileID: "2",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, ""),
			},
			wantResponse: `{"error_code":400,"error_message":"invalid user id"}`,
			wantStatus:   400,
			prepare:      func(a args) {},
		},
		// NEGATIVE || Fail due service layer error
		{
			name: "Fail_due_to_internal_server_error",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"error_code":500,"error_message":"unable to send email"}`,
			wantStatus:   http.StatusInternalServerError,
			prepare: func(a args) {
				s.service.On("SendUserInvitation", mock.Anything, TestUserID, TestProfileID).Return(errors.New("error sending invitation")).Once()
			},
		},
		// NEGATIVE || Empty body with valid IDs
		{
			name: "Empty_body_with_valid_ids",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"data":{"message":"Invitation sent successfully to employee"}}`,
			wantStatus:   http.StatusOK,
			prepare: func(a args) {
				s.service.On("SendUserInvitation", mock.Anything, TestUserID, TestProfileID).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)
			r := httptest.NewRequest(http.MethodPost, "/employee_invite", nil)
			r = r.WithContext(tt.args.ctx)
			r = mux.SetURLVars(r, map[string]string{"profile_id": "1"})
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.SendUserInvitation(tt.args.ctx, s.service))
			handler.ServeHTTP(w, r)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("Expected %d but got %d", tt.wantStatus, w.Result().StatusCode)
			}

			if w.Body.String() != tt.wantResponse {
				t.Errorf("Expected response body %s but got %s", tt.wantResponse, w.Body.String())
			}
		})
	}
}

func (s *HandlerTestSuite) TestSendAdminInvitation() {
	t := s.T()

	type args struct {
		userID    string
		profileID string
		ctx       context.Context
	}

	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantResponse string
		prepare      func(args)
	}{
		// POSITIVE || Successfully update invitation
		{
			name: "Successfully update invitation",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"data":{"message":"Profile Completed Successfully"}}`,
			wantStatus:   http.StatusOK,
			prepare: func(a args) {
				s.service.On("UpdateInvitation", mock.Anything, TestUserID, TestProfileID).Return(nil).Once()
			},
		},
		// // NEGATIVE || Fail due to invalid user id
		{
			name: "Fail_due_to_invalid_user_id",
			args: args{
				userID:    "",
				profileID: "2",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, ""),
			},
			wantResponse: `{"error_code":400,"error_message":"invalid user id"}`,
			wantStatus:   http.StatusBadRequest,
			prepare:      func(a args) {},
		},
		// // NEGATIVE || Fail due to invalid profile id
		{
			name: "Fail_due_to_invalid_profile_id",
			args: args{
				userID:    "1",
				profileID: "",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"error_code":400,"error_message":"invalid request data"}`,
			wantStatus:   http.StatusBadRequest,
			prepare:      func(a args) {},
		},
		// // NEGATIVE || Fail due to missing context value
		{
			name: "Fail_due_to_missing_context_value",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.Background(),
			},
			wantResponse: `{"error_code":400,"error_message":"invalid user id"}`,
			wantStatus:   http.StatusBadRequest,
			prepare:      func(a args) {},
		},
		// // NEGATIVE || Fail due to service layer error
		{
			name: "Fail_due_to_internal_server_error",
			args: args{
				userID:    "1",
				profileID: "1",
				ctx:       context.WithValue(context.Background(), constants.UserIDKey, 1.0),
			},
			wantResponse: `{"error_code":500,"error_message":"unable to send email"}`,
			wantStatus:   http.StatusInternalServerError,
			prepare: func(a args) {
				s.service.On("UpdateInvitation", mock.Anything, TestUserID, TestProfileID).Return(errors.New("error sending invitation")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)
			r := httptest.NewRequest(http.MethodPost, "/admin_invite", nil)
			r = r.WithContext(tt.args.ctx)
			r = mux.SetURLVars(r, map[string]string{"profile_id": tt.args.profileID})
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(handler.SendAdminInvitation(tt.args.ctx, s.service))
			handler.ServeHTTP(w, r)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("Expected %d but got %d", tt.wantStatus, w.Result().StatusCode)
			}

			if w.Body.String() != tt.wantResponse {
				t.Errorf("Expected response body %s but got %s", tt.wantResponse, w.Body.String())
			}
		})
	}
}
