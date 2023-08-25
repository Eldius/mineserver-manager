package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// GetFileName returns file name from URL
func GetFileName(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return filepath.Base(p.Path)
}

func DownloadFile(timeout time.Duration, u, destFile string) error {
	c := http.Client{Timeout: timeout}
	res, err := c.Get(u)
	if err != nil {
		err = fmt.Errorf("getting version info: %w", err)
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

// ValidateFileIntegrity validates down
func ValidateFileIntegrity(file, signature string) error {
	log := logger.GetLogger()
	in, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("opening source file: %w", err)
		return err
	}

	hash := sha1.New()
	if _, err := io.Copy(hash, in); err != nil {
		err = fmt.Errorf("reading source file content: %w", err)
		return err
	}
	sum := hash.Sum(make([]byte, 0))

	fileSignature := fmt.Sprintf("%x", sum)
	log.With("calculated", fileSignature, "original", signature).Info("FileChecksumValidation")
	if fileSignature != signature {
		return errors.New("file sign validation error")
	}

	return nil
}

// UnpackTarGZ unpacks .tar.gz file to destDir
func UnpackTarGZ(file, destDir string) error {
	log := logger.GetLogger().With("action", "unpack", "file", file, "dest", destDir)
	log.Info("File '%s' unpacked to '%s'", file, destDir)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		err = fmt.Errorf("creating installation base path: %w", err)
		log.With("error", err).Error("Failed to create destination directory")
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
	}()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		err = fmt.Errorf("reading package file: %w", err)
		log.With("error", err).Error("Failed open package file")
		return err
	}

	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			err = fmt.Errorf("reading next file in package: %w", err)
			log.With("error", err).Error("Failed to get next file in package")
			return err
		}

		name := header.Name

		dest := filepath.Join(destDir, name)
		switch header.Typeflag {
		case tar.TypeDir:
			fmt.Printf("  -> Creating folder %s\n", dest)
			if err := os.MkdirAll(dest, os.ModePerm); err != nil {
				err = fmt.Errorf("failed to create subdirectory for new file: %w", err)
				log.Error("Failed to create folder %s: %v", dest, err)
				return err
			}
		case tar.TypeReg:
			fmt.Printf("  -> Creating file %s\n", dest)
			dir := filepath.Dir(dest)
			if _, err := os.Stat(dir); err != nil {
				_ = os.MkdirAll(dir, os.ModePerm)
			}

			f, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				err = fmt.Errorf("creating dest file for '%s': %w", name, err)
				log.Error("Failed to create file %s: %v\n", dest, err)
			}
			_, err = io.Copy(f, tarReader)
			if err != nil {
				err = fmt.Errorf("writing content to destination file: %w", err)
				log.Error("Failed to write to file %s: %v\n", dest, err)
				return err
			}
		default:
			slog.Warn("%s : %c %s %s",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
	return nil
}
