package errors

import "errors"

// Error Variables for replacing the errors wherever required
var (
	ErrInvalidBody            = errors.New("invalid request body")
	ErrParameterMissing       = errors.New("parameter missing")
	ErrInvalidID              = errors.New("invalid id")
	ErrEmptyPayload           = errors.New("empty payload array")
	ErrSecretKey              = errors.New("secret key not found")
	ErrEmailNotFound          = errors.New("email not found")
	ErrDecodeRequest          = errors.New("unable to decode request")
	ErrGoogleRequest          = errors.New("unable to send request to google")
	ErrDecodeResponse         = errors.New("unable to decode response")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrGenerateToken          = errors.New("unable to generate token")
	ErrEmptyAccessToken       = errors.New("empty access token")
	ErrAuthToken              = errors.New("unauthorized access")
	ErrAuthHeader             = errors.New("invalid authorization header")
	ErrSigningMethod          = errors.New("unexpected signing method")
	ErrInvalispecsken         = errors.New("invalid token")
	ErrUserID                 = errors.New("error in parsing userID from claims")
	ErrTokenEmpty             = errors.New("token string is empty")
	ErrTokenExpirationHours   = errors.New("TOKEN_EXPIRATION_HOURS is not set")
	ErrVerifyToken            = errors.New("error in verifying token")
	ErrFailespecsFetch        = errors.New("failed to fetch data")
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrInvalidEnv             = errors.New("no env variable found")
	ErrFailedToDelete         = errors.New("failed to delete")
	ErrNoData                 = errors.New("no data found")
	ErrFailedToUpdate         = errors.New("failed to update status")
	ErrComponentNotSuppoerted = errors.New("component name not supported")
	ErrUnableToSendEmail      = errors.New("unable to send email")
	ErrFailedToGet            = errors.New("failed to get data")
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
	ErrInvalidConfig    = errors.New("invalid configuration for database connection")
	ErrMisMatchParams   = errors.New("mismatch in number of records for component")
)
