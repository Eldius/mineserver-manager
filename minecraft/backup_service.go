package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
	"log/slog"
	"path/filepath"
	"time"
)

type BackupService interface {
	// Backup creates a new backup from instance
	Backup(ctx context.Context, instancePath, backupDestFolder string) (string, error)
	// Restore restores a backup file to instance
	Restore(ctx context.Context, instancePath, backupDestFolder string) error
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
		err = fmt.Errorf("parsing backupDestPath: %w", err)
		return "", err
	}
	destFile := filepath.Join(
		backupDestPath,
		fmt.Sprintf(
			"%s_%s_backup.zip",
			filepath.Base(instancePath),
			time.Now().Format("2006-01-02_15-04-05"),
		))

	return destFile, utils.PackFiles(ctx, instancePath, destFile)
}

func (s *backupService) Restore(ctx context.Context, instancePath, backupDestFolder string) error {
	return nil
}
