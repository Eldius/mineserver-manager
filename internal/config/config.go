package config

import (
	"github.com/spf13/viper"
	"time"
)

func GetMinecraftDownloadTimeout() time.Duration {
	return viper.GetDuration(minecraftDownloadTimeoutPropKey)
}

func GetMinecraftApiTimeout() time.Duration {
	return viper.GetDuration(minecraftApiTimeoutPropKey)
}
