package utils

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func PackFiles(ctx context.Context, src, dest string) error {
	log := logger.GetLogger().
		With(
			slog.String("action", "pack"),
			slog.String("src", src),
			slog.String("dest", dest),
		)
	log.Debug("Starting to pack files")

	return pack(ctx, src, dest)
}

func pack(_ context.Context, src, dest string) error {
	log := logger.GetLogger().
		With(
			slog.String("action", "pack"),
			slog.String("src", src),
			slog.String("dest", dest))

	f, err := os.Create(dest)
	if err != nil {
		err = fmt.Errorf("opening file to backup (%s): %w", src, err)
		return err
	}
	w := zip.NewWriter(f)

	defer func() {
		_ = w.Flush()
		_ = w.Close()
	}()

	javaFolder := filepath.Join(src, "java")
	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		log = log.With(
			slog.String("path", path),
			slog.String("name", info.Name()),
			slog.Bool("is_dir", info.IsDir()),
		)

		log.Debug("start processing file")

		if err != nil {
			return err
		}
		if path == src {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".log") || strings.HasSuffix(info.Name(), ".log.gz") {
			return nil
		}

		if strings.HasPrefix(path, javaFolder) {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(src, "libraries")) {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(src, "versions")) {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(src, "crash-reports")) {
			return nil
		}

		if strings.HasPrefix(path, javaFolder) {
			return nil
		}

		if info.Name() == "server.pid" {
			return nil
		}

		fmt.Printf("- file name: %s\n", info.Name())
		fmt.Printf("  file path: %s\n", path)
		out, err := w.Create(strings.TrimPrefix(path, src))
		if err != nil {
			err = fmt.Errorf("creating file to backup (%s): %w", path, err)
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			err = fmt.Errorf("opening file to backup (%s): %w", path, err)
			return err
		}

		log.Debug("copying file to zip")
		if _, err := io.Copy(out, in); err != nil {
			err = fmt.Errorf("copying file to backup (%s): %w", path, err)
			return err
		}

		return nil
	}); err != nil {
		err = fmt.Errorf("listing instance files: %w", err)
		return err
	}

	return nil
}
