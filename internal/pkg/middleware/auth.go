package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

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
			log.Fatalf("Error in verifying token: %v", err)
			ErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		fmt.Println("Claims: ", claims)
		userID, ok := claims["userID"]
		if !ok {
			ErrorResponse(w, http.StatusUnauthorized, errors.ErrInvalidToken)
			zap.S().Error(errors.ErrUserID)
			return
		}
		fmt.Println("type of userId is : ", reflect.TypeOf(userID))

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
