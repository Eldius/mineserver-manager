package installer

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/httpx"
	"path/filepath"
	"time"
)

type Downloader interface {
	DownloadServer(ctx context.Context, url, sha1, dest string) (string, error)
}

type vanillaDownloader struct {
	timeout time.Duration
}

func NewDownloader(timeout time.Duration) Downloader {
	return &vanillaDownloader{
		timeout: timeout,
	}
}

func (d *vanillaDownloader) DownloadServer(ctx context.Context, url, sha1, dest string) (string, error) {
	destFile := filepath.Join(dest, httpx.GetFileName(url))
	if err := httpx.DownloadFile(ctx, d.timeout, url, destFile); err != nil {
		return "", &DownloadError{
			URL:  url,
			Dest: destFile,
			Err:  fmt.Errorf("downloading server file: %w", err),
		}
	}

	if err := httpx.ValidateFileIntegrity(ctx, destFile, sha1); err != nil {
		return "", &DownloadError{
			URL:  url,
			Dest: destFile,
			Err:  fmt.Errorf("validating server file integrity: %w", err),
		}
	}

	return destFile, nil
}
