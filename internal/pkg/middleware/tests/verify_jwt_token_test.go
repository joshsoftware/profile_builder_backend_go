package tests

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
)

var (
	userID = 123
	Email  = "test@example.com"
)

func init() {
	os.Setenv("SECRET_KEY", "test_secret_key")
}
func createTokenString(secretKey []byte, claims jwt.MapClaims, signingMethod jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(secretKey)
}

// Predefined function to generate token strings for test cases
func generateTokenString(secretKey []byte, claims jwt.MapClaims, signingMethod jwt.SigningMethod) string {
	tokenString, _ := createTokenString(secretKey, claims, signingMethod)
	return tokenString
}

func TestVerifyJwtToken(t *testing.T) {
	test_secret_key := []byte(os.Getenv("SECRET_KEY"))
	fixedTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return fixedTime })
	if err != nil {
		t.Fatal(err)
	}
	defer patch.Unpatch()

	tests := []struct {
		name             string
		tokenString      string
		expectedError    error
		expectedResponse jwt.MapClaims
	}{
		{
			name: "success",
			tokenString: generateTokenString(test_secret_key, jwt.MapClaims{
				"userID": userID,
				"email":  Email,
				"exp":    time.Now().Add(time.Hour).Unix(),
			}, jwt.SigningMethodHS256),
			expectedError:    nil,
			expectedResponse: jwt.MapClaims{"userID": float64(userID), "email": Email, "exp": float64(fixedTime.Add(time.Hour).Unix())},
		},
		{
			name:          "empty token string",
			tokenString:   "",
			expectedError: errors.New("token string is empty"),
		},
		{
			name: "error_token_expired",
			tokenString: generateTokenString(test_secret_key, jwt.MapClaims{
				"userID": userID,
				"email":  Email,
				"exp":    time.Now().Add(-time.Hour).Unix(),
			}, jwt.SigningMethodHS256),
			expectedError: errors.New("Token is expired"),
		},
		{
			name:          "invalid token format",
			tokenString:   "invalid_token_format",
			expectedError: errors.New("token contains an invalid number of segments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := middleware.VerifyJWTToken(tt.tokenString)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, claims)
			}
		})
	}
}
