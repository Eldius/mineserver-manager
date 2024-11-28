package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"golang.org/x/term"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrChecksumValidationFailed = errors.New("file sign validation error")
	ErrCouldNotOpenFile         = errors.New("opening source file")
	ErrCouldNotReadFile         = errors.New("reading source file content")
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

func WarnOnError[T any](obj T, err error) T {
	if err != nil {
		logger.GetLogger().With("error", err).Warn("WarnOnError")
	}
	return obj
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

// UnpackTarGZ unpacks .tar.gz file to destDir
func UnpackTarGZ(ctx context.Context, file, destDir string) error {
	log := logger.GetLogger().With(
		slog.String("action", "unpack"),
		slog.String("file", file),
		slog.String("dest", destDir),
	)

	log.Info("Starting to unpack file")
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		err = fmt.Errorf("creating installation base path: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create destination directory")
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
		log.With("error", err).ErrorContext(ctx, "Failed open package file")
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
			log.With("error", err).ErrorContext(ctx, "Failed to get next file in package")
			return err
		}

		name := header.Name

		dest := filepath.Join(destDir, name)
		switch header.Typeflag {
		case tar.TypeDir:
			fmt.Printf("  -> Creating folder %s\n", dest)
			if err := os.MkdirAll(dest, os.ModePerm); err != nil {
				err = fmt.Errorf("failed to create subdirectory for new file: %w", err)
				log.ErrorContext(ctx, "Failed to create folder %s: %v", dest, err)
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
				log.ErrorContext(ctx, "Failed to create file %s: %v\n", dest, err)
			}
			_, err = io.Copy(f, tarReader)
			if err != nil {
				err = fmt.Errorf("writing content to destination file: %w", err)
				log.ErrorContext(ctx, "Failed to write to file %s: %v\n", dest, err)
				return err
			}
		default:
			slog.WarnContext(ctx, "%s : %c %s %s",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
	return nil
}

// HTTPClient returns a new HTTP client
func HTTPClient(t time.Duration) http.Client {
	return http.Client{
		Timeout: t,
	}
}

// PasswordPrompt prompts user for a password
// (hiding it's value from console)
func PasswordPrompt() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		err = fmt.Errorf("password prompt: %w", err)
		return "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(password), nil
}

// ExpandPath expands tilde '~' character for home
func ExpandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// AbsolutePath returns the absolute path
func AbsolutePath(path string) (string, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return "", fmt.Errorf("expanded path: %w", err)
	}
	return filepath.Abs(path)
}

func shaHash(content []byte) string {
	hash := sha256.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}

// ShaHash calculates hash from io.Reader content
func ShaHash(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		err = fmt.Errorf("reading file to backup (%s): %w", r, err)
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
