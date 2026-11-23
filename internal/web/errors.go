package web

import "errors"

var (
	// ErrPathOutsideHome is returned when a requested path would escape
	// the configured API home folder.
	ErrPathOutsideHome = errors.New("path outside home folder")

	// ErrAuthNotConfigured is returned when the API auth password is not
	// set and therefore no login can be performed.
	ErrAuthNotConfigured = errors.New("api auth not configured")
)
