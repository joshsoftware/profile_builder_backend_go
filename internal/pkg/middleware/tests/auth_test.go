package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/undefinedlabs/go-mpatch"
)

var (
	userID1 = 1
	Email1  = "test@example.com"
)

func newRequest(method, url, token string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestAuthMiddleware(t *testing.T) {

	tests := []struct {
		Name               string
		Request            *http.Request
		MockVerifyJWTToken func(string) (jwt.MapClaims, error)
		ExpectedStatusCode int
	}{
		{
			Name:    "Success",
			Request: newRequest("GET", "/", "valispecsken"),
			MockVerifyJWTToken: func(tokenString string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"userID": int64(123), "Email": Email1}, nil
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:    "No Authorization Header",
			Request: newRequest("GET", "/", ""),
			MockVerifyJWTToken: func(tokenString string) (jwt.MapClaims, error) {
				return nil, errors.ErrAuthToken
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:    "Invalid Authorization Header",
			Request: newRequest("GET", "/", "invalispecsken"),
			MockVerifyJWTToken: func(tokenString string) (jwt.MapClaims, error) {
				return nil, errors.ErrAuthHeader
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:    "Error Verifying Token",
			Request: newRequest("GET", "/", "errorToken"),
			MockVerifyJWTToken: func(tokenString string) (jwt.MapClaims, error) {
				return nil, errors.ErrVerifyToken
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:    "Invalid Token",
			Request: newRequest("GET", "/", "invalispecsken"),
			MockVerifyJWTToken: func(tokenString string) (jwt.MapClaims, error) {
				// Return an error instead of a valid claim
				return nil, errors.ErrInvalidToken
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			patch, _ := mpatch.PatchMethod(middleware.VerifyJWTToken, tt.MockVerifyJWTToken)
			defer patch.Unpatch()
			rr := httptest.NewRecorder()
			handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			handler.ServeHTTP(rr, tt.Request)
			if rr.Code != tt.ExpectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.ExpectedStatusCode)
			}
		})
	}
}
