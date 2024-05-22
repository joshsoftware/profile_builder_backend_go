package errors

import "errors"

// Error Variables for replacing the errors wherever required
var (
	ErrInvalidBody      = errors.New("invalid request body")
	ErrParameterMissing = errors.New("parameter missing")
	ErrInvalidID        = errors.New("invalid id")
	ErrEmptyPayload     = errors.New("empty payload array")
)

// Profile Related variables
var (
	ErrInvalidFormat = errors.New("invalid request format")
)

// DB Related variables
var (
	ErrConnectionFailed = errors.New("error connecting to database")
	ErrDuplicateKey     = errors.New("record already exists")
	ErrInvalidProfile   = errors.New("invalid profile id")
)
