package cmd

import (
	"context"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/minecraft"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"time"
)

// installCmd installs a minecraft server instance
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a Minecraft server instance",
	Long:  `Installs a Minecraft server instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.GetLogger().With("headless", installHeadless).Info("debugging")
		ctx := context.Background()
		client := minecraft.NewInstallService(minecraft.WithTimeout(
			config.GetMinecraftApiTimeout()),
			minecraft.WithDownloadTimeout(config.GetMinecraftDownloadTimeout()),
		)

		if installJustListVersions {
			if err := minecraft.ListVersions(ctx); err != nil {
				err = fmt.Errorf("listing available versions: %w", err)
				slog.ErrorContext(ctx, "failed to list available versions: %v", err)
				panic(err)
			}
		} else {
			if err := client.Install(ctx,
				serverconfig.WithVersion(installServerVersion),
				serverconfig.ToDestinationFolder(installDestinationFolder),
				serverconfig.WithHeadlessConfig(installHeadless),
			); err != nil {
				err = fmt.Errorf("installing server: %w", err)
				slog.ErrorContext(ctx, "failed to install server: %v", err)
				panic(err)
			}
		}
	},
}

var (
	installServerVersion     string
	installDestinationFolder string
	installHeadless          bool
	installJustListVersions  bool
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installServerVersion, "version", "latest", "Version of Java Edition server to install")
	installCmd.Flags().StringVar(&installDestinationFolder, "dest", ".", "Installation root directory")
	installCmd.Flags().BoolVar(&installHeadless, "headless", false, "Installation root directory")
	installCmd.Flags().BoolVar(&installJustListVersions, "list", false, "Lists available versions to install")

	installCmd.Flags().Duration("download-timeout", 300*time.Second, "Download timeout configuration (defaults to 300s/5m)")
	if err := viper.BindPFlag(config.AppInstallDownloadTimeoutPropKey, installCmd.Flags().Lookup("download-timeout")); err != nil {
		err = fmt.Errorf("binding artifact download timeout property: %w", err)
		panic(err)
	}
}
