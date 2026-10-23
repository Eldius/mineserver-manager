package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// backupRestoreCmd represents the restore command
var backupRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore instance backup",
	Long:  `Restore instance backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBackupRestore(context.Background(), backupRestoreOpts)
	},
}

var (
	backupRestoreOpts struct {
		fromFile string
		toFolder string
	}
)

func init() {
	backupCmd.AddCommand(backupRestoreCmd)

	backupRestoreCmd.Flags().StringVar(&backupRestoreOpts.fromFile, "backup-file", "", "Backup file to be restored")
	backupRestoreCmd.Flags().StringVar(&backupRestoreOpts.toFolder, "instance-folder", ".", "Installation root directory (defaults to current directory)")
	if err := backupRestoreCmd.MarkFlagRequired("backup-file"); err != nil {
		panic(err)
	}
}
