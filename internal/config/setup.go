package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*
Setup sets up app configuration
*/
func Setup(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mqtt-listener-go" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".mineserver"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.SetConfigType("yml")
	}

	SetDefaults()
	viper.SetEnvPrefix("mineserver")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}

/*
SetDefaults sets default configuration values
*/
func SetDefaults() {
	viper.SetDefault("curseforge.timeout", "30s")
	viper.SetDefault("app.devmode", false)
	viper.SetDefault("app.request.log", false)
	viper.SetDefault("app.install.path", "./.tmp")
}
