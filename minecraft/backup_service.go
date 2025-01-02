package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	bkpTimestampFormat = "2006-01-02_15-04-05"
)

type BackupService interface {
	// Backup creates a new backup from instance
	Backup(ctx context.Context, instancePath, backupDestFolder string) (string, error)
	// Restore restores a backup file to instance
	Restore(ctx context.Context, instancePath, backupFile string) error
	// RotateBackupFiles limits max backup files stored
	RolloverBackupFiles(ctx context.Context, backupDestFolder string) error
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
			time.Now().Format(bkpTimestampFormat),
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

func (s *backupService) RolloverBackupFiles(ctx context.Context, backupDestFolder string) error {
	log := slog.With(
		slog.String("backup_folder", backupDestFolder),
	)
	stat, err := os.Stat(backupDestFolder)
	if err != nil {
		return fmt.Errorf("getting backup dir info: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("backup dir is not a directory: %s", backupDestFolder)
	}
	bkpFiles, err := mapBackupFiles(ctx, backupDestFolder)
	if err != nil {
		return fmt.Errorf("getting backup files: %w", err)
	}

	for k, _ := range bkpFiles {
		ts, err := time.Parse(k, bkpTimestampFormat)
		if err != nil {
			return fmt.Errorf("parsing backup timestamp: %w", err)
		}
		log.With("bkp_file_timestamp", ts.Format("2006-01-02_15-04-05")).DebugContext(ctx, "backup files")
	}
	log.With("bkp_files", bkpFiles).DebugContext(ctx, "backup files")
	return nil
}

func mapBackupFiles(ctx context.Context, backupDestFolder string) (map[string]os.DirEntry, error) {
	filesMap := make(map[string]os.DirEntry)

	entries, err := os.ReadDir(backupDestFolder)
	if err != nil {
		return filesMap, fmt.Errorf("reading backup dir: %w", err)
	}

	rgxp, err := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}-[0-9]{2}-[0-9]{2}_backup.zip")
	if err != nil {
		return filesMap, fmt.Errorf("compile regexp: %w", err)
	}

	for _, entry := range entries {
		log := slog.With("entry_name", entry.Name())
		str := rgxp.FindString(entry.Name())
		log.With(
			slog.String("find_str", str),
		).DebugContext(ctx, "parsing backup file")
		if len(str) > 0 {
			filesMap[strings.TrimSuffix(str, "_backup.zip")] = entry
		}
	}
	return filesMap, nil
}
