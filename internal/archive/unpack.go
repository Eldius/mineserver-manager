package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func Unpack(_ context.Context, instancePath, backupFile string) error {
	r, err := zip.OpenReader(backupFile)
	if err != nil {
		return fmt.Errorf("opening backup file: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		outFile := filepath.Join(instancePath, f.Name)

		slog.With(
			slog.String("instance_path", instancePath),
			slog.String("backup_file", f.Name),
			slog.String("current_file", f.Name),
			slog.String("dest_file", outFile),
		).Debug("UnpackingFile")

		if err := os.MkdirAll(filepath.Dir(outFile), os.ModePerm); err != nil {
			return fmt.Errorf("mkdirall: %w", err)
		}

		out, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}

		in, err := f.Open()
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}

		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
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
