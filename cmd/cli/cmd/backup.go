package cmd

import (
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backs up a server",
	Long:  `Backs up a server.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var (
	backupOpts struct {
		instance string
	}
)

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&backupOpts.instance, "instance", ".", "Installation root directory (defaults to current directory)")

}
