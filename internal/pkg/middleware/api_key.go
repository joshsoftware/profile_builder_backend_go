package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// APIKeyMiddleware authenticates server-to-server requests using a pre-shared API key.
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedKey := os.Getenv(constants.APIKeyEnvVar)
		if expectedKey == "" {
			// Misconfiguration: the server-side key is not set — deny all requests.
			zap.S().Error("API key environment variable is not configured: ", constants.APIKeyEnvVar)
			ErrorResponse(w, http.StatusInternalServerError, errors.ErrInvalidEnv)
			return
		}

		incomingKey := r.Header.Get(constants.APIKeyHeader)
		if incomingKey == "" {
			zap.S().Warn("Request rejected: missing API key header")
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrAPIKeyMissing)
			return
		}

		if subtle.ConstantTimeCompare([]byte(expectedKey), []byte(incomingKey)) != 1 {
			zap.S().Warn("Request rejected: invalid API key")
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrAPIKeyInvalid)
			return
		}

		next.ServeHTTP(w, r)
	})
}
