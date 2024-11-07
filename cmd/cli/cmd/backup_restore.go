package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// backupRestoreCmd represents the restore command
var backupRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore instance backup",
	Long:  `Restore instance backup.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restore called")
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

	backupRestoreCmd.Flags().StringVar(&backupRestoreOpts.toFolder, "instance-folder", ".", "Installation root directory (defaults to current directory)")
	backupRestoreCmd.Flags().StringVar(&backupRestoreOpts.fromFile, "backup-folder", ".backups", "Backup file destination folder (defaults to .backups on current directory)")
}
