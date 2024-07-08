package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// AuthMiddleware used for Authentication.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrAuthToken)
			zap.S().Error(errors.ErrAuthToken)
			return
		}

		splitToken := strings.Split(token, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrAuthHeader)
			zap.S().Error(errors.ErrAuthHeader)
			return
		}

		tokenString := splitToken[1]
		claims, err := VerifyJWTToken(tokenString)
		if err != nil {
			zap.S().Error("Error in verifying token: %v", err)
			ErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		userID, ok := claims["userID"]
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalispecsken)
			zap.S().Error(errors.ErrUserID)
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
