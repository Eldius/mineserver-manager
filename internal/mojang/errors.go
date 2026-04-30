package mojang

import (
	"errors"
	"fmt"
)

var (
	// ErrRequestFailed is returned when a mojang API request fails
	ErrRequestFailed = errors.New("mojang API request failed")
	// ErrDecodeFailed is returned when decoding a mojang API response fails
	ErrDecodeFailed = errors.New("mojang API response decoding failed")
)

// APIError provides context for mojang API failures
type APIError struct {
	URL  string
	Verb string
	Err  error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("mojang API %s request to %s failed: %v", e.Verb, e.URL, e.Err)
}

func (e *APIError) Unwrap() error {
	return e.Err
}
