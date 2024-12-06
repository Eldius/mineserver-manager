package cmd

import (
	"context"
	"errors"
	"github.com/eldius/mineserver-manager/minecraft"

	"github.com/spf13/cobra"
)

// backupRestoreCmd represents the restore command
var backupRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore instance backup",
	Long:  `Restore instance backup.`,
	Run: func(cmd *cobra.Command, args []string) {
		if backupRestoreOpts.fromFile == "" {
			panic(errors.New("invalid input file"))
		}
		ctx := context.Background()
		if err := minecraft.NewBackupService().Restore(ctx, backupRestoreOpts.toFolder, backupRestoreOpts.fromFile); err != nil {
			panic(err)
		}
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
	backupRestoreCmd.Flags().StringVar(&backupRestoreOpts.fromFile, "backup-file", "", "Backup file destination folder (defaults to .backups on current directory)")
}
