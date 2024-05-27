package errors

import "errors"

var (
	ErrEmptyPayload      = errors.New("empty payload array")
	ErrSectetKeyNotFound = errors.New("secret key not found")
)

// DB Related variables
var (
	ErrConnectionFailed = errors.New("error connecting to database")
	ErrDuplicateKey     = errors.New("record already exists")
	ErrInvalidProfile   = errors.New("invalid profile id")
	ErrNoRecordFound    = errors.New("no record found")
)
