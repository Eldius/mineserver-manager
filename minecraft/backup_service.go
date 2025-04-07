package minecraft

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/utils"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	bkpTimestampFormat = "2006-01-02_15-04-05"
)

type BackupService interface {
	// Backup creates a new backup from instance
	Backup(ctx context.Context, instancePath, backupDestFolder string) (*BackupInfo, error)
	// Restore restores a backup file to instance
	Restore(ctx context.Context, instancePath, backupFile string) error
	// RolloverBackupFiles limits max backup files stored
	RolloverBackupFiles(ctx context.Context, backupDestFolder, backupName string, maxBkpFiles int) error
}

type backupService struct {
}

func NewBackupService() BackupService {
	return &backupService{}
}

func (s *backupService) Backup(ctx context.Context, instancePath, backupDestPath string) (*BackupInfo, error) {

	log := slog.With(
		slog.String("instance_path", instancePath),
		slog.String("backup_dest_folder", backupDestPath),
	)

	log.InfoContext(ctx, "starting backup process")

	instancePath, err := utils.AbsolutePath(instancePath)
	if err != nil {
		err = fmt.Errorf("parsing to absolute Path: %w", err)
		return nil, err
	}

	backupDestPath, err = utils.AbsolutePath(backupDestPath)
	if err != nil {
		return nil, fmt.Errorf("parsing backupDestPath: %w", err)
	}
	instanceName := filepath.Base(instancePath)
	ts := time.Now()
	destFile := filepath.Join(
		backupDestPath,
		fmt.Sprintf(
			"%s_%s_backup.zip",
			instanceName,
			ts.Format(bkpTimestampFormat),
		))

	if err := utils.PackFiles(ctx, instancePath, destFile); err != nil {
		return nil, fmt.Errorf("writing backup file: %w", err)
	}

	return &BackupInfo{
		Timestamp: ts,
		Name:      instanceName,
		Path:      destFile,
	}, nil
}

func (s *backupService) Restore(ctx context.Context, instancePath, backupFile string) error {

	if err := os.MkdirAll(instancePath, os.ModePerm); err != nil {
		return fmt.Errorf("creating backup dir: %w", err)
	}

	return utils.Unpack(ctx, instancePath, backupFile)
}

func (s *backupService) RolloverBackupFiles(ctx context.Context, backupDestFolder, backupName string, maxBkpFiles int) error {
	log := slog.With(
		slog.String("backup_folder", backupDestFolder),
		slog.Int("max_bkp_files", maxBkpFiles),
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

	bkpList := bkpFiles[backupName]
	deleteAllBeforeIdx := len(bkpList) - maxBkpFiles

	log = log.With(slog.Int("delete_all_before_idx", deleteAllBeforeIdx), slog.Int("bkp_count", len(bkpList)))

	for i, b := range bkpList[:deleteAllBeforeIdx] {
		l := log.With(
			slog.String("bkp_path", b.Path),
			slog.String("bkp_name", b.Name),
			slog.Int("bkp_idx", i),
		)
		l.DebugContext(ctx, "deleting backup file")
		if err := os.Remove(b.Path); err != nil {
			err := fmt.Errorf("deleting backup file: %w", err)
			l.With("error", err).ErrorContext(ctx, "deleting backup file")
			return err
		}
	}
	return nil
}

func mapBackupFiles(ctx context.Context, backupDestFolder string) (backupsMapping, error) {
	filesMap := make(backupsMapping)

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
			bkpName := strings.TrimSuffix(entry.Name(), "_"+str)
			tsStr := strings.TrimSuffix(str, "_backup.zip")

			log = log.With("error", err, "ts_str", tsStr, "bkp_name", bkpName, "ts_str", tsStr)

			var bkpList []BackupInfo

			if l, ok := filesMap[bkpName]; ok {
				bkpList = l
			}

			ts, err := time.Parse(bkpTimestampFormat, tsStr)
			if err != nil {
				log.With("error", err, "ts_str", tsStr, "bkp_name", bkpName, "ts_str", tsStr).
					WarnContext(ctx, "backup file Timestamp parsing failed")
				continue
			}
			filesMap[bkpName] = append(bkpList, BackupInfo{
				Timestamp: ts,
				Name:      bkpName,
				Path:      filepath.Join(backupDestFolder, entry.Name()),
			})

		}
	}

	for k := range filesMap {
		sort.Slice(filesMap[k], func(i, j int) bool {
			return filesMap[k][i].Timestamp.Before(filesMap[k][j].Timestamp)
		})
	}

	return filesMap, nil
}

type BackupInfo struct {
	Timestamp time.Time
	Name      string
	Path      string
}

type backupsMapping map[string]backupList

type backupList []BackupInfo

func (i backupList) olderFile() *BackupInfo {
	if len(i) == 0 {
		return nil
	}

	var older *BackupInfo

	for _, bkpInfo := range i {
		if older == nil {
			older = &bkpInfo
			continue
		}
		if bkpInfo.Timestamp.Before(older.Timestamp) {
			older = &bkpInfo
		}
	}

	return older
}
