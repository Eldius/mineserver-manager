package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
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

	if err := os.MkdirAll(instancePath, os.ModePerm); err != nil {
		return fmt.Errorf("creating backup dir: %w", err)
	}

	return nil
}
