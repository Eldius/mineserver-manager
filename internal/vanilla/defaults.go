package vanilla

import (
	"embed"
	"fmt"
	"github.com/eldius/properties"
)

var (
	//go:embed all:default_values
	defaultConfigFiles embed.FS
	// DefaultEulaValue default eula.txt content
	DefaultEulaValue = &Eula{Eula: true}
)

func GetDefaultServerProperties() (*ServerProperties, error) {
	var resp ServerProperties
	in, err := defaultConfigFiles.Open("default_values/server.properties")
	if err != nil {
		err = fmt.Errorf("reading default server.properties values: %w", err)
		return nil, err
	}
	if err := properties.NewDecoder(in).Decode(&resp); err != nil {
		err = fmt.Errorf("reading default server.properties values: %w", err)
		return nil, err
	}
	return &resp, nil
}

func GetDefaultScriptParams() *StartupParams {
	return &StartupParams{
		Xmx:           "1g",
		Xms:           "1g",
		LogConfigFile: "",
	}
}
