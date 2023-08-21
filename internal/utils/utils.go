package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// GetFileName returns file name from URL
func GetFileName(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return filepath.Base(p.Path)
}

func DownloadFile(c *http.Client, u, destFile string) error {
	res, err := c.Get(u)
	if err != nil {
		err = fmt.Errorf("getting version info: %w", err)
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server dest file: %w", err)
		return err
	}

	if _, err := io.Copy(f, res.Body); err != nil {
		err = fmt.Errorf("copying server file to dest file: %w", err)
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("status code not success (was %d)", res.StatusCode)
	}

	return nil
}

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
