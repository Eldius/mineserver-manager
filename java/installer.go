package java

import (
	"context"
	"fmt"
	utils "github.com/eldius/mineserver-manager/internal/utils"
	"os"
	"path/filepath"
	"time"
)

var (
	// JavaVersions is a map of Java Runtime download links
	JavaVersions = map[int]map[string]map[string]string{
		21: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-21.0.4-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-21.0.4-linux-aarch64.tar.gz",
			},
		},
		17: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-17.0.7-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-17.0.7-linux-aarch64.tar.gz",
			},
		},
		11: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-11.0.19-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-11.0.19-linux-aarch64.tar.gz",
			},
		},
	}
)

// DownloadJDK downloads JVM package
func DownloadJDK(ctx context.Context, v int, arch, osName string, timeout time.Duration) (string, error) {
	u := JavaVersions[v][osName][arch]
	tempDir, err := os.MkdirTemp(os.TempDir(), "mine-installer-*")
	if err != nil {
		err = fmt.Errorf("creating temp folder to save java runtime (osName: %s/arch: %s/v: %d): %w", osName, arch, v, err)
		return "", err
	}

	dest := filepath.Join(tempDir, utils.GetFileName(u))
	if err := utils.DownloadFile(ctx, timeout, u, dest); err != nil {
		err = fmt.Errorf("downloading java runtime: %w", err)
		return "", err
	}

	return dest, nil
}
