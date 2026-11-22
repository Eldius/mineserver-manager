package cmd

import (
	"context"
	"fmt"

	"github.com/eldius/cli-telemetry-helper/pkg/setup"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "mineserver-manager",
	Short:              "A simple CLI tool to manage Minecraft server installations",
	Long:               `A simple CLI tool to manage Minecraft server installations.`,
	PersistentPreRunE:  setup.PersistentPreRunE,
	PersistentPostRunE: setup.PersistentPostRunE,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) error {
	return setup.CobraCmdExecute(ctx, rootCmd)
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
