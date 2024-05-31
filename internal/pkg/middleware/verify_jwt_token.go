package middleware

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

func VerifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrSigningMethod
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		log.Fatalf("Error in parsing token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Fatalf("Token is invalid")
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatalf("Error in parsing claims")
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
