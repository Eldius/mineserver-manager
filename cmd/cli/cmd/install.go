package cmd

import (
	"fmt"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/vanilla"
	"github.com/spf13/cobra"
	"log/slog"
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
		client := vanilla.NewInstallService(vanilla.WithTimeout(
			config.GetMinecraftApiTimeout()),
			vanilla.WithDownloadTimeout(config.GetMinecraftDownloadTimeout()),
		)

		if err := client.Install(vanilla.WithVersion(installServerVersion), vanilla.ToDestinationFolder(installDestinationFolder)); err != nil {
			err = fmt.Errorf("installing server: %w", err)
			slog.Error("failed to install server: %v", err)
			panic(err)
		}
	},
}

var (
	installServerVersion     string
	installDestinationFolder string
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installServerVersion, "version", "latest", "Version of Java Edition server to install")
	installCmd.Flags().StringVar(&installDestinationFolder, "dest", "/tmp/installations", "Installation root directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
