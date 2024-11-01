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

func init() {
	rootCmd.AddCommand(backupRestoreCmd)
}
