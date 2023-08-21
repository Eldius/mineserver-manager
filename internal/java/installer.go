package java

import (
	"fmt"
	utils "github.com/eldius/mineserver-manager/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// JavaVersions is a map of Java Runtime download links
	JavaVersions = map[int]map[string]map[string]string{
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
func DownloadJDK(v int, arch, osName string) (string, error) {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	u := JavaVersions[v][osName][arch]
	tempDir, err := os.MkdirTemp(os.TempDir(), "mine-installer-*")
	if err != nil {
		err = fmt.Errorf("creating temp folder to save java runtime (osName: %s/arch: %s/v: %d): %w", osName, arch, v, err)
		return "", err
	}

	dest := filepath.Join(tempDir, utils.GetFileName(u))
	if err := utils.DownloadFile(c, u, dest); err != nil {
		err = fmt.Errorf("downloading java runtime: %w", err)
		return "", err
	}

	return dest, nil
}
