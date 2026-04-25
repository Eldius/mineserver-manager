package java

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/internal/utils"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// PackageVersions is a map of Java Runtime download links
	PackageVersions = map[int]map[string]map[string]string{
		8: {
			"linux": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10_openj9-0.26.0/OpenJDK8U-jdk_x64_linux_openj9_8u292b10_openj9-0.26.0.tar.gz",
				"arm64": "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10_openj9-0.26.0/OpenJDK8U-jdk_aarch64_linux_openj9_8u292b10_openj9-0.26.0.tar.gz",
			},
		},
		16: {
			"linux": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9_openj9-0.26.0/OpenJDK16U-jdk_x64_linux_openj9_16.0.1_9_openj9-0.26.0.tar.gz",
				"arm64": "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9_openj9-0.26.0/OpenJDK16U-jdk_aarch64_linux_openj9_16.0.1_9_openj9-0.26.0.tar.gz",
			},
		},
		17: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-debugsymbols-17.0.18-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-debugsymbols-17.0.18-linux-aarch64.tar.gz",
			},
		},
		21: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-21.0.4-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-21.0.4-linux-aarch64.tar.gz",
			},
		},
		25: {
			"linux": {
				"amd64": "https://aka.ms/download-jdk/microsoft-jdk-debugsymbols-25.0.2-linux-x64.tar.gz",
				"arm64": "https://aka.ms/download-jdk/microsoft-jdk-debugsymbols-25.0.2-linux-aarch64.tar.gz",
			},
		},
	}
)

// Download downloads JVM package
func Download(ctx context.Context, v int, arch, osName string, timeout time.Duration) (string, error) {
	u := PackageVersions[v][osName][arch]
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

// Install downloads and unpack JDK to a destination folder
func Install(ctx context.Context, dest string, v int, arch, osName string, timeout time.Duration) (string, error) {
	log := logger.GetLogger().With(slog.String("action", "install_jdk"), slog.Int("jdk_version", v))

	jdkPackage, err := Download(ctx, v, arch, osName, timeout)
	if err != nil {
		err = fmt.Errorf("downloading java runtime to install: %w", err)
		return "", err
	}
	defer func() {
		_ = os.RemoveAll(filepath.Dir(jdkPackage))
	}()

	if err = utils.UnpackTarGZ(ctx, jdkPackage, dest); err != nil {
		err = fmt.Errorf("unpacking jdk package: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to unpack JDK package")
		return "", err
	}

	//return dest, nil
	jdkUnpacked, err := findJDKUnpackedFolder(dest)
	if err != nil {
		err = fmt.Errorf("finding unpacked jdk root folder: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to unpack JDK package")
		return "", err
	}

	jdkBasePath := filepath.Join(dest, "jdk")
	log.With(
		slog.String("dest", dest),
		slog.String("jdk_unpacked", jdkUnpacked),
		slog.String("jdk_package", jdkPackage),
		slog.String("jdk_new_folder", jdkBasePath),
	).Info("Found JDK folder")

	if err := os.Rename(jdkUnpacked, jdkBasePath); err != nil {
		err = fmt.Errorf("renaming unpacked jdk root folder: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to unpack JDK package")
		return "", err
	}
	return jdkUnpacked, nil
}

func findJDKUnpackedFolder(root string) (string, error) {
	fmt.Println("looking for JDK folder")
	entries, err := os.ReadDir(root)
	if err != nil {
		err = fmt.Errorf("reading jdk root folder (%s): %w", root, err)
		return "", err
	}

	for _, entry := range entries {
		fmt.Printf(" - %s (is dir: %v)\n", entry.Name(), entry.IsDir())
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "jdk") {
			fmt.Println("found", entry.Name())
			return filepath.Join(root, entry.Name()), nil
		}
	}

	return "", err
}
