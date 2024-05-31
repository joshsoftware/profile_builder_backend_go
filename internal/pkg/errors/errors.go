package errors

import "errors"

// Error Variables for replacing the errors wherever required
var (
	ErrInvalidBody      = errors.New("invalid request body")
	ErrParameterMissing = errors.New("parameter missing")
	ErrInvalidID        = errors.New("invalid id")
	ErrEmptyPayload     = errors.New("empty payload array")
	ErrSecretKey        = errors.New("secret key not found")
	ErrEmailNotFound    = errors.New("email not found")
	ErrDecodeRequest    = errors.New("Unable to decode request")
	ErrGoogleRequest    = errors.New("Unable to send request to google")
	ErrDecodeResponse   = errors.New("Unable to decode response")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrGenerateToken    = errors.New("Unable to generate token")
	ErrEmptyAccessToken = errors.New("empty access token")
	ErrAuthToken        = errors.New("Authorization token not found")
	ErrAuthHeader       = errors.New("invalid authorization header")
	ErrSigningMethod    = errors.New("unexpected signing method")
	ErrInvalidToken     = errors.New("invalid token")
	ErrUserID           = errors.New("Error in parsing userID from claims")
)

// Profile Related variables
var (
	ErrInvalidFormat          = errors.New("invalid request format")
	ErrInvalidRequestData     = errors.New("invalid request data")
	ErrHTTPRequestFailed      = errors.New("failed to perform HTTP request")
	ErrReadResponseBodyFailed = errors.New("failed to read response body")
)

// DB Related variables
var (
	ErrConnectionFailed = errors.New("error connecting to database")
	ErrDuplicateKey     = errors.New("record already exists")
	ErrInvalidProfile   = errors.New("invalid profile id")
	ErrNoRecordFound    = errors.New("no record found")
)
