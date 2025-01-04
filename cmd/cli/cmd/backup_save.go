package cmd

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/minecraft"

	"github.com/spf13/cobra"
)

// backupSaveCmd represents the save command
var backupSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a backup from instance",
	Long:  `Save a backup from instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		s := minecraft.NewBackupService()
		backupFile, err := s.Backup(ctx, backupSaveOpts.instance, backupSaveOpts.destFolder)
		if err != nil {
			fmt.Printf("Failed to make a backup: %v\n", err)
			return err
		}

		if backupSaveOpts.maxBackupFiles > 0 {
			if err := s.RolloverBackupFiles(ctx, backupSaveOpts.destFolder, backupFile.Name, backupSaveOpts.maxBackupFiles); err != nil {
				fmt.Printf("Failed to make a backup rollover: %v\n", err)
				return err
			}
		}
		fmt.Printf("Backup completed to '%s'!\n", backupFile)
		return nil
	},
}

var (
	backupSaveOpts struct {
		instance       string
		destFolder     string
		maxBackupFiles int
	}
)

func init() {
	backupCmd.AddCommand(backupSaveCmd)

	backupSaveCmd.Flags().StringVar(&backupSaveOpts.instance, "instance-folder", ".", "Installation root directory (defaults to current directory)")
	backupSaveCmd.Flags().StringVar(&backupSaveOpts.destFolder, "backup-folder", ".backups", "Backup file destination folder (defaults to .backups on current directory)")
	backupSaveCmd.Flags().IntVar(&backupSaveOpts.maxBackupFiles, "max-backup-files", 0, "Max number of backup files to be stored (defaults to 0 - disabled)")
}
