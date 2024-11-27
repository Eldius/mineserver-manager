package utils

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
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

	var hashes bytes.Buffer
	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		log := log.With(
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
		packedFileName := strings.TrimPrefix(path, src)
		out, err := w.Create(packedFileName)
		if err != nil {
			err = fmt.Errorf("creating file to backup (%s): %w", path, err)
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			err = fmt.Errorf("opening file to backup (%s): %w", path, err)
			return err
		}

		b, err := io.ReadAll(in)
		if err != nil {
			err = fmt.Errorf("reading file to backup (%s): %w", path, err)
			return err
		}

		log.Debug("copying file to zip")
		if _, err := out.Write(b); err != nil {
			err = fmt.Errorf("writing file to zip (%s): %w", path, err)
			return err
		}
		hash := sha256.New()
		if _, err := hash.Write(b); err != nil {
			err = fmt.Errorf("calculating file hash (%s): %w", path, err)
			return err
		}
		hashes.WriteString(fmt.Sprintf("%s  %x\n", packedFileName, hash.Sum(nil)))

		return nil
	}); err != nil {
		err = fmt.Errorf("listing instance files: %w", err)
		return err
	}

	hf, err := w.Create("backup.sha256")
	if err != nil {
		err = fmt.Errorf("creating file to backup (%s): %w", src, err)
		return err
	}
	if _, err := hf.Write(hashes.Bytes()); err != nil {
		err = fmt.Errorf("writing file to backup (%s): %w", src, err)
		return err
	}

	return nil
}
