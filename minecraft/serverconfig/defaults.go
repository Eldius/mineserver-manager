package serverconfig

import (
	"embed"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig/generators"
)

var (
	//go:embed all:default_values
	defaultConfigFiles embed.FS
	// DefaultEulaValue default eula.txt content
	DefaultEulaValue = &generators.Eula{Eula: true}
)

func GetDefaultRuntimeParams() *generators.RuntimeGenerator {
	return &generators.RuntimeGenerator{
		Xmx:           "1g",
		Xms:           "1g",
		LogConfigFile: true,
	}
}
