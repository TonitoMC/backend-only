package errors

import "errors"

// This package contains errors as they will be presented to the user, they're imported in the
// services layer to be able to return them and 'bubble them up' with wrapping up to the
// HTTP layer. Then they're translated into error codes & logged, any error with a different
// identification is Code 500 'Internal Server Error'

var (
	ErrSeriesNotFound      = errors.New("series not found")
	ErrSeriesAlreadyExists = errors.New("series already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrSeriesConflict      = errors.New("series conflict")
)
