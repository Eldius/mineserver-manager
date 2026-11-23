package config

const (
	AppInstallDownloadTimeoutPropKey = "install.download.timeout"
	AppMinecraftAPITimeoutPropKey    = "minecraft.api.timeout"

	AppHomePathPropKey    = "app.home.path"
	AppInstallPathPropKey = "app.install.path"
	AppRequestLogPropKey  = "app.request.log"
	AppDebugModePropKey   = "app.debugmode"

	ApiHomePathPropKey    = "api.home.path"
	ApiAuthPasswordPropKey = "api.auth.password"
	ApiAuthTokenTTLPropKey = "api.auth.token.ttl"

	AppHomeDefaultValue    = "~/.mineserver"
	ApiHomeDefaultValue    = "~/.mineserver-api"
	ApiAuthTokenTTLDefault = "24h"

	VersionsFileName = "versions.json"

	CliAppName = "mineserver-cli"
	ApiAppName = "mineserver-api"
)
