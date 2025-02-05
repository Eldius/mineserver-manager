package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// purpurInstallCmd represents the install command
var purpurInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs Purpur server",
	Long:  `Installs Purpur server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(purpurInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// purpurInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purpurInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
