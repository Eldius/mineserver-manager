/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/eldius/mineserver-manager/internal/vanilla"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := vanilla.NewClient(2 * time.Second)
		if err := client.InstallWithConfig(vanilla.WithVersion(installServerVersion)); err != nil {
			err = fmt.Errorf("installing server: %w", err)
			log.Fatalf("failed to install server: %v", err)
		}

	},
}

var (
	installServerVersion string
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installServerVersion, "version", "latest", "Version of Java Edition server to install")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
