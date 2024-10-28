package cmd

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/minecraft"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backs up a server",
	Long:  `Backs up a server.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		backupFile, err := minecraft.NewBackupService().Backup(ctx, backupOpts.instance, backupOpts.destFolder)
		if err != nil {
			fmt.Printf("Failed to make a backup: %v\n", err)
		}

		fmt.Printf("Backup completed to %s!\n", backupFile)
	},
}

var (
	backupOpts struct {
		instance   string
		destFolder string
	}
)

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&backupOpts.instance, "from", ".", "Installation root directory (defaults to current directory)")
	backupCmd.Flags().StringVar(&backupOpts.destFolder, "to", ".backups", "Backup file destination folder (defaults to .backups on current directory)")
}
