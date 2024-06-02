package middleware

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

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
		fmt.Println("Error in parsing token in verify token : ", err)
		return nil, err
	}

	if !token.Valid {
		fmt.Println("Token is invalid")
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		fmt.Println("Error in parsing claims")
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
