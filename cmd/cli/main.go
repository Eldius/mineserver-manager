package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eldius/cli-telemetry-helper/pkg/setup"
	"github.com/eldius/mineserver-manager/cmd/cli/cmd"
	"github.com/eldius/mineserver-manager/internal/config"
)

func main() {
	os.Exit(func() int {
		ctx := context.Background()
		svc := setup.Service{
			Name:    "minepacker-cli",
			Version: "0.0.0a",
		}
		ini := func() error {
			return cmd.Execute(ctx)
		}
		if err := setup.MainWrapper(
			ctx,
			svc,
			ini,
			setup.WithConfigFileName("config"),
			setup.WithConfigFileType("yaml"),
			setup.WithConfigFileLookupDirectories("."),
			setup.WithEnvPrefix("cli_telemetry"),
			setup.WithDefaultConfigs(map[string]any{
				config.AppMinecraftAPITimeoutPropKey:    "10s",
				config.AppInstallDownloadTimeoutPropKey: "300s",
				config.AppDebugModePropKey:              false,
				config.AppRequestLogPropKey:             false,
				config.AppInstallPathPropKey:            "./.tmp",
				config.AppHomePathPropKey:               config.AppHomeDefaultValue,
			}),
		); err != nil {
			fmt.Printf("setup telemetry failed: %v\n", err)
			return 1
		}
		return 0
	}())
}
