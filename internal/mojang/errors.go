package mojang

import (
	"errors"
	"fmt"
)

var (
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
