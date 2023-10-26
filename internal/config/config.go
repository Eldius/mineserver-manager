package config

import (
	"github.com/spf13/viper"
	"time"
)

func GetMinecraftDownloadTimeout() time.Duration {
	return viper.GetDuration(AppInstallDownloadTimeoutPropKey)
}

func GetMinecraftApiTimeout() time.Duration {
	return viper.GetDuration(AppMinecraftAPITimeoutPropKey)
}

func GetAppHomePath() string {
	return viper.GetString(AppHomePathPropKey)
}
