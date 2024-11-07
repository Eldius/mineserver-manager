package cmd

import (
	"fmt"
	initCfg "github.com/eldius/initial-config-go/configs"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mineserver-manager",
	Short: "A simple CLI tool to manage Minecraft server installations",
	Long:  `A simple CLI tool to manage Minecraft server installations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initCfg.InitSetup("", initCfg.WithConfigFileToBeUsed(cfgFile),
			initCfg.WithDefaultCfgFileLocations("/.mineserver", "."),
			initCfg.WithDefaultCfgFileName("config"),
			initCfg.WithDefaultValues(map[string]any{
				config.AppMinecraftAPITimeoutPropKey:    "10s",
				config.AppInstallDownloadTimeoutPropKey: "300s",
				config.AppDebugModePropKey:              false,
				config.AppRequestLogPropKey:             false,
				config.AppInstallPathPropKey:            "./.tmp",
				config.AppHomePathPropKey:               config.AppHomeDefaultValue,
				initCfg.LogLevelKey:                     initCfg.LogLevelDEBUG,
				initCfg.LogFormatKey:                    initCfg.LogFormatJSON,
				initCfg.LogOutputToStdoutKey:            true,
			}))
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (defaults to $HOME/.mineserver/config.yaml)")

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug log")
	if err := viper.BindPFlag(config.AppDebugModePropKey, rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		err = fmt.Errorf("binding debug mode property key: %w", err)
		panic(err)
	}

	rootCmd.PersistentFlags().String("home", config.AppHomeDefaultValue, "Define app's home folder (defaults to $HOME/.mineserver)")
	if err := viper.BindPFlag(config.AppHomePathPropKey, rootCmd.PersistentFlags().Lookup("home")); err != nil {
		err = fmt.Errorf("binding debug mode property key: %w", err)
		panic(err)
	}
}
