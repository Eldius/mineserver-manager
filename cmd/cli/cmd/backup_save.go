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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		backupFile, err := minecraft.NewBackupService().Backup(ctx, backupSaveOpts.instance, backupSaveOpts.destFolder)
		if err != nil {
			fmt.Printf("Failed to make a backup: %v\n", err)
		}

		fmt.Printf("Backup completed to '%s'!\n", backupFile)
	},
}

var (
	backupSaveOpts struct {
		instance   string
		destFolder string
	}
)

func init() {
	backupCmd.AddCommand(backupSaveCmd)

	backupSaveCmd.Flags().StringVar(&backupSaveOpts.instance, "instance-folder", ".", "Installation root directory (defaults to current directory)")
	backupSaveCmd.Flags().StringVar(&backupSaveOpts.destFolder, "backup-folder", ".backups", "Backup file destination folder (defaults to .backups on current directory)")
}
