package cmd

import (
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Instance backup management",
	Long:  `Instance backup management.`,
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
