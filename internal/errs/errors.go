package errs

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")
	// ErrInvalidInput is returned when provided input is invalid
	ErrInvalidInput = errors.New("invalid input")
	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")
	// ErrPermissionDenied is returned when an operation is not permitted
	ErrPermissionDenied = errors.New("permission denied")
	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")
	// ErrInternal is returned for unexpected internal errors
	ErrInternal = errors.New("internal error")
)
