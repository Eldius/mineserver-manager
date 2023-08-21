package vanilla

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

var (
	//go:embed all:templates
	templateFiles embed.FS

	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseFS(templateFiles, "templates/**"))
}

type Eula struct {
	Eula bool `properties:"eula"`
}

type StartupParams struct {
	Xmx           string
	Xms           string
	LogConfigFile string
	Headless      bool
}

func (s StartupParams) ToScript() (string, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "start.sh", s); err != nil {
		err = fmt.Errorf("generating start script: %w", err)
		return "", err
	}

	return b.String(), nil
}

type InstallConfig struct {
	Start       *StartupParams
	SrvProps    *ServerProperties
	Dest        string
	VersionName string
	v           *VersionInfoResponse
}

type InstallCfg func(*InstallConfig) *InstallConfig

func WithVersion(v string) InstallCfg {
	return func(c *InstallConfig) *InstallConfig {
		c.VersionName = v
		return c
	}
}
