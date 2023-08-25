package serverconfig

import (
	"embed"
)

var (
	//go:embed all:default_values
	defaultConfigFiles embed.FS
	// DefaultEulaValue default eula.txt content
	DefaultEulaValue = &Eula{Eula: true}
)

func GetDefaultScriptParams() *StartupParams {
	return &StartupParams{
		Xmx:           "1g",
		Xms:           "1g",
		LogConfigFile: "",
	}
}
