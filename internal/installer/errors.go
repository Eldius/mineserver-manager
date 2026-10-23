package installer

import (
	"errors"
	"fmt"
)

var (
	// ErrFlavorNotFound is returned when a requested flavor is not registered
	ErrFlavorNotFound = errors.New("flavor not found")
	// ErrVersionNotFound is returned when a requested version is not found for a flavor
	ErrVersionNotFound = errors.New("version not found")
)

// DownloadError provides context for download failures
type DownloadError struct {
	URL  string
	Dest string
	Err  error
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("failed to download from %s to %s: %v", e.URL, e.Dest, e.Err)
}

func (e *DownloadError) Unwrap() error {
	return e.Err
}
