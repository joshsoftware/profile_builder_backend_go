package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
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

		helpers.WhiteListMutext.Lock()
		_, ok := helpers.TokenList[tokenString]
		helpers.WhiteListMutext.Unlock()

		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrAuthToken)
			zap.S().Error(errors.ErrTokenNotFound)
			return
		}

		userID, ok := claims["userID"]
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalispecsken)
			zap.S().Error(errors.ErrUserID)
			return
		}

		role, ok := claims["role"]
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalispecsken)
			zap.S().Error(errors.ErrUserRole)
			return
		}

		reqProfileID, ok := claims["profileID"]
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalispecsken)
			zap.S().Error(errors.ErrProfileID)
			return
		}

		requestedProfileID := helpers.ConvertFloatToInt(reqProfileID)

		if !helpers.ProfileIDNotRequiredPath(r) {
			profileID, err := helpers.GetProfileId(r)
			if err != nil {
				ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalidProfile)
				zap.S().Error(errors.ErrProfileID)
				return
			}

			if requestedProfileID != profileID && role != constants.Admin {
				ErrorResponse(w, http.StatusUnauthorized, errors.ErrAuthToken)
				zap.S().Error(errors.ErrProfileID)
				return
			}
		}

		ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
		ctx = context.WithValue(ctx, constants.UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware used for Authorization.
func RoleMiddleware(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(constants.UserRoleKey).(string)
			if !ok {
				ErrorResponse(w, http.StatusUnauthorized, errors.ErrAuthToken)
				zap.S().Error(errors.ErrAuthToken)
				return
			}
			for _, role := range roles {
				if role == userRole {
					next.ServeHTTP(w, r)
					return
				}
			}
			ErrorResponse(w, http.StatusForbidden, errors.ErrAuthToken)
			zap.S().Error(errors.ErrAuthToken)
		})
	}
}
