package minecraft

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type BackupService interface {
	// Backup creates a new backup from instance
	Backup(ctx context.Context, instancePath, backupDestFolder string) (string, error)
	// Restore restores a backup file to instance
	Restore(ctx context.Context, instancePath, backupFile string) error
}

type backupService struct {
}

func NewBackupService() BackupService {
	return &backupService{}
}

func (s *backupService) Backup(ctx context.Context, instancePath, backupDestPath string) (string, error) {

	log := slog.With(
		slog.String("instance_path", instancePath),
		slog.String("backup_dest_folder", backupDestPath),
	)

	log.InfoContext(ctx, "starting backup process")

	instancePath, err := utils.AbsolutePath(instancePath)
	if err != nil {
		err = fmt.Errorf("parsing to absolute path: %w", err)
		return "", err
	}

	backupDestPath, err = utils.AbsolutePath(backupDestPath)
	if err != nil {
		return "", fmt.Errorf("parsing backupDestPath: %w", err)
	}
	destFile := filepath.Join(
		backupDestPath,
		fmt.Sprintf(
			"%s_%s_backup.zip",
			filepath.Base(instancePath),
			time.Now().Format("2006-01-02_15-04-05"),
		))

	if err := utils.PackFiles(ctx, instancePath, destFile); err != nil {
		return "", fmt.Errorf("writing backup file: %w", err)
	}

	return destFile, nil
}

func (s *backupService) Restore(_ context.Context, instancePath, backupFile string) error {

	_ = os.MkdirAll(instancePath, os.ModePerm)

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

		_ = os.MkdirAll(filepath.Dir(outFile), os.ModePerm)

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
