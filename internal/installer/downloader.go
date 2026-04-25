package installer

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/utils"
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
	destFile := filepath.Join(dest, utils.GetFileName(url))
	if err := utils.DownloadFile(ctx, d.timeout, url, destFile); err != nil {
		return "", fmt.Errorf("downloading server file: %w", err)
	}

	if err := utils.ValidateFileIntegrity(ctx, destFile, sha1); err != nil {
		return "", err
	}

	return destFile, nil
}
