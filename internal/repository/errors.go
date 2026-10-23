package repository

import "errors"

var (
	// ErrNotFound is returned when an instance is not found in the repository
	ErrNotFound = errors.New("instance not found")
	// ErrSaveFailed is returned when persisting an instance fails
	ErrSaveFailed = errors.New("saving instance failed")
	// ErrDeleteFailed is returned when deleting an instance fails
	ErrDeleteFailed = errors.New("deleting instance failed")
)
