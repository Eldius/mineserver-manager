// Package httpx provides HTTP helpers: client construction, file downloads
// and download integrity validation.
package httpx

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrChecksumValidationFailed = errors.New("file sign validation error")
	ErrCouldNotOpenFile         = errors.New("opening source file")
	ErrCouldNotReadFile         = errors.New("reading source file content")
	ErrDownloadFailed           = errors.New("download failed")
)

// GetFileName returns file name from URL
func GetFileName(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return filepath.Base(p.Path)
}

// DownloadFile downloads a file
func DownloadFile(ctx context.Context, timeout time.Duration, u, destFile string) error {
	c := HTTPClient(timeout)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		err = fmt.Errorf("creating versions query request: %w", err)
		return err
	}
	res, err := c.Do(r)
	if err != nil {
		err = fmt.Errorf("downloading file from %s: %w", u, err)
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	destFolder := filepath.Dir(destFile)
	if _, err := os.Stat(destFolder); err != nil {
		_ = os.MkdirAll(destFolder, os.ModePerm)
	}

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server dest file: %w", err)
		return err
	}

	if _, err := io.Copy(f, res.Body); err != nil {
		err = fmt.Errorf("copying downloaded server file to dest file: %w", err)
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("%w: status code was %d", ErrDownloadFailed, res.StatusCode)
	}

	return nil
}

// ValidateFileIntegrity validates down
func ValidateFileIntegrity(ctx context.Context, file, signature string) error {
	log := logger.GetLogger()
	in, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("%s: %w", ErrCouldNotOpenFile, err)
		return err
	}

	hash := sha1.New()
	if _, err := io.Copy(hash, in); err != nil {
		err = fmt.Errorf("%s: %w", ErrCouldNotReadFile, err)
		return err
	}
	sum := hash.Sum(make([]byte, 0))

	fileSignature := fmt.Sprintf("%x", sum)
	log.With("calculated", fileSignature, "original", signature).InfoContext(ctx, "FileChecksumValidation")
	if fileSignature != signature {
		return fmt.Errorf("file sign validation error (%s): %w", fmt.Sprintf("calculated: %s => expected: %s)", fileSignature, signature), ErrChecksumValidationFailed)
	}

	return nil
}

// HTTPClient returns a new HTTP client
func HTTPClient(t time.Duration) http.Client {
	return http.Client{
		Timeout: t,
	}
}
