package config

import (
	"time"

	"github.com/eldius/cli-telemetry-helper/pkg/setup"
)

func GetMinecraftDownloadTimeout() time.Duration {
	return setup.GetDuration(AppInstallDownloadTimeoutPropKey)
}

func GetMinecraftApiTimeout() time.Duration {
	return setup.GetDuration(AppMinecraftAPITimeoutPropKey)
}

func GetAppHomePath() string {
	return setup.GetString(AppHomePathPropKey)
}

// GetApiHomePath returns the root folder the API confines all
// filesystem access (installs, backups, temp files) to.
func GetApiHomePath() string {
	return setup.GetString(AppHomePathPropKey)
}

// GetApiAuthPassword returns the single admin password used by the API
// login endpoint. Empty means auth is not configured and the API
// refuses logins.
func GetApiAuthPassword() string {
	return setup.GetString(ApiAuthPasswordPropKey)
}

// GetApiAuthTokenTTL returns how long an API auth token stays valid.
func GetApiAuthTokenTTL() time.Duration {
	return setup.GetDuration(ApiAuthTokenTTLPropKey)
}
