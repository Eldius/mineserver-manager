package serverconfig

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

type RuntimeParams struct {
	Xmx           string
	Xms           string
	LogConfigFile bool
	Headless      bool
}

func (s RuntimeParams) ToScript() (string, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "start.sh", s); err != nil {
		err = fmt.Errorf("generating start script: %w", err)
		return "", err
	}

	return b.String(), nil
}

func (s RuntimeParams) LoggingConfiguration(installPath string) (string, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "log4j2.xml", installPath); err != nil {
		err = fmt.Errorf("generating logging configuration file: %w", err)
		return "", err
	}

	return b.String(), nil
}
