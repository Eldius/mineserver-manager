/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eldius/cli-telemetry-helper/pkg/setup"
	"github.com/eldius/mineserver-manager/internal/buildinfo"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/web"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mineserver-api",
	Short: "Web UI and JSON API to manage Minecraft server instances",
	Long: `Starts the mineserver manager API server.

It serves an embedded web UI and a JSON API to install, list,
start/stop and back up Minecraft server instances. All filesystem
access is confined to the configured api home folder (api.home.path),
and the API requires authentication (api.auth.password).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := web.StartApi(rootCmdOpts.port); err != nil {
			cmd.Println("Error starting web server:", err)
			slog.With("error", err).ErrorContext(cmd.Context(), "error starting web server")
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
func Execute() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("recovered from a panic:", r)
			os.Exit(1)
		}
	}()

	os.Exit(
		setup.Run(
			context.Background(),
			setup.Service{
				Name:    config.ApiAppName,
				Version: buildinfo.Version,
			},
			rootCmd,
			setup.WithConfigFileLookupDirectories(".", config.AppHomeDefaultValue),
			setup.WithEnvPrefix(config.ApiAppEnvPrefix),
			setup.WithDefaultConfigs(map[string]any{
				config.AppMinecraftAPITimeoutPropKey:    "10s",
				config.AppInstallDownloadTimeoutPropKey: "300s",
				config.AppDebugModePropKey:              false,
				config.AppRequestLogPropKey:             false,
				config.AppInstallPathPropKey:            filepath.Join(config.AppHomeDefaultValue, ".tmp"),
				config.AppHomePathPropKey:               config.AppHomeDefaultValue,
				config.ApiAuthPasswordPropKey:           "",
				config.ApiAuthTokenTTLPropKey:           config.ApiAuthTokenTTLDefault,
			}),
		))
}

var (
	rootCmdOpts struct {
		port int
	}
)

func init() {
	rootCmd.Flags().IntVarP(&rootCmdOpts.port, "port", "p", 8080, "port to listen on")
}
