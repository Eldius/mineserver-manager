package config

import (
	"embed"
	"fmt"
	"github.com/eldius/properties"
	"io/fs"
)

var (
	//go:embed all:default_values
	defaultConfigFiles embed.FS
	// DefaultEulaValue default eula.txt content
	DefaultEulaValue = &Eula{Eula: true}
)

func GetDefaultConfigFile(f string) (fs.File, error) {
	return defaultConfigFiles.Open("default_values/server.properties")
}

// DefaultServerProperties returns the default server.properties representation
func DefaultServerProperties() (*ServerProperties, error) {
	var resp ServerProperties
	in, err := GetDefaultConfigFile("default_values/server.properties")
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
