package middleware

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// VerifyJWTToken used to verify user given token
func VerifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.ErrTokenEmpty
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrSigningMethod
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		zap.S().Error("Error in parsing token: %v", err)
		return nil, errors.ErrInvalispecsken
	}

	if !token.Valid {
		zap.S().Error("Token is not valid")
		return nil, errors.ErrInvalispecsken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		zap.S().Error("Error in parsing claims")
		return nil, errors.ErrInvalispecsken
	}

	return claims, nil
}
