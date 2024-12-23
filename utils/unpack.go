package utils

import (
	"archive/zip"
	"context"
	"fmt"
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
