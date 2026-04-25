package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

// backupSaveCmd represents the save command
var backupSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a backup from instance",
	Long:  `Save a backup from instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBackupSave(context.Background(), backupSaveOpts)
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
