package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/minecraft"
)

func runBackupSave(ctx context.Context, opts struct {
	instance       string
	destFolder     string
	maxBackupFiles int
}) error {
	s := minecraft.NewBackupService()
	backupFile, err := s.Backup(ctx, opts.instance, opts.destFolder)
	if err != nil {
		return fmt.Errorf("failed to make a backup: %w", err)
	}

	if opts.maxBackupFiles > 0 {
		if err := s.RolloverBackupFiles(ctx, opts.destFolder, backupFile.Name, opts.maxBackupFiles); err != nil {
			return fmt.Errorf("failed to make a backup rollover: %w", err)
		}
	}
	fmt.Printf("Backup completed to '%s'!\n", backupFile.Path)
	return nil
}

func runBackupRestore(ctx context.Context, opts struct {
	fromFile string
	toFolder string
}) error {
	if err := minecraft.NewBackupService().Restore(ctx, opts.toFolder, opts.fromFile); err != nil {
		if errors.Is(err, minecraft.ErrBackupFileNotFound) {
			return fmt.Errorf("backup file not found: %s", opts.fromFile)
		}
		return fmt.Errorf("restoring backup: %w", err)
	}
	return nil
}
