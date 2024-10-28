package minecraft

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/minecraft/model"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/eldius/properties"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BackupService interface {
	Backup(ctx context.Context, instancePath, backupDestFolder string) (string, error)
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

	instancePath, err := filepath.Abs(instancePath)
	if err != nil {
		err = fmt.Errorf("parsing instance instancePath: %w", err)
		return "", err
	}

	backupDestPath, err = filepath.Abs(backupDestPath)
	if err != nil {
		err = fmt.Errorf("parsing backupDestPath: %w", err)
		return "", err
	}

	versionsFilePath := filepath.Join(instancePath, config.VersionsFileName)
	stat, err := os.Stat(versionsFilePath)
	if err != nil {
		err = fmt.Errorf("checking if versions file exists: %w", err)
		return "", err
	}
	if stat.IsDir() {
		err = fmt.Errorf("versions file is a directory")
		return "", err
	}

	f, err := os.Open(versionsFilePath)
	if err != nil {
		err = fmt.Errorf("opening versions file: %w", err)
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	var versionsInfo model.VersionsInfo
	if err := json.NewDecoder(f).Decode(&versionsInfo); err != nil {
		err = fmt.Errorf("decoding versions file: %w", err)
		return "", err
	}

	log = slog.With("versions_info", versionsInfo)
	log.InfoContext(ctx, "found versions file")

	propsFilePath := filepath.Join(instancePath, "server.properties")

	propsFile, err := os.Open(propsFilePath)
	if err != nil {
		err = fmt.Errorf("opening properties file: %w", err)
		return "", err
	}

	var props serverconfig.ServerProperties
	if err := properties.NewDecoder(propsFile).Decode(&props); err != nil {
		err = fmt.Errorf("parsing properties file: %w", err)
		return "", err
	}

	javaFolder := filepath.Join(instancePath, "java")
	var filesToBackup []fileToBackup
	if err := filepath.Walk(instancePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == instancePath {
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

		if strings.HasPrefix(path, filepath.Join(instancePath, "libraries")) {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(instancePath, "versions")) {
			return nil
		}

		if strings.HasPrefix(path, filepath.Join(instancePath, "crash-reports")) {
			return nil
		}

		if strings.HasPrefix(path, javaFolder) {
			return nil
		}

		if info.Name() == "server.pid" {
			return nil
		}

		fmt.Printf("file name: %s\n", info.Name())
		fmt.Printf("file path: %s\n", path)
		filesToBackup = append(filesToBackup, fileToBackup{
			src: path,
			dst: strings.TrimPrefix(path, instancePath),
		})

		return nil
	}); err != nil {
		err = fmt.Errorf("listing instance files: %w", err)
		return "", err
	}

	destFile := filepath.Join(
		backupDestPath,
		fmt.Sprintf(
			"%s_%s_backup.zip",
			filepath.Base(instancePath),
			time.Now().Format("2006-01-02_15-04-05"),
		))
	if err := createZipFile(ctx, destFile, filesToBackup); err != nil {
		err = fmt.Errorf("creating backup zip file: %w", err)
		return "", err
	}

	return destFile, nil
}

func createZipFile(_ context.Context, destFile string, files []fileToBackup) error {
	_ = os.MkdirAll(filepath.Dir(destFile), 0755)
	f, err := os.Create(destFile)
	if err != nil {
		err = fmt.Errorf("creating backup destination file: %w", err)
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	w := zip.NewWriter(f)
	defer func() {
		_ = w.Close()
	}()

	for _, file := range files {
		if err := copyFile(file, w); err != nil {
			err = fmt.Errorf("copying file: %w", err)
			return err
		}
	}

	return nil
}

func copyFile(file fileToBackup, w *zip.Writer) error {
	in, err := os.Open(file.src)
	if err != nil {
		err = fmt.Errorf("opening file to backup (%s): %w", file.src, err)
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	fmt.Printf(" -> adding file to archive: %s..\n", file.dst)
	dstFileWriter, err := w.Create(file.dst)
	if err != nil {
		err = fmt.Errorf("creating file to backup (%s): %w", file.dst, err)
		return err
	}

	_, err = io.Copy(dstFileWriter, in)
	if err != nil {
		err = fmt.Errorf("copying file to backup (%s): %w", file.dst, err)
		return err
	}
	return nil
}

type fileToBackup struct {
	src string
	dst string
}
