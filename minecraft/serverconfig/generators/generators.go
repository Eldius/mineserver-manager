package generators

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

const (
	StartScriptFileName = "start.sh"
	StopScriptFileName  = "stop.sh"
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

type RuntimeGenerator struct {
	Xmx           string
	Xms           string
	LogConfigFile bool
	Headless      bool
}

func StopScript() (string, error) {
	f, err := templateFiles.Open(filepath.Join("templates", StopScriptFileName))
	if err != nil {
		err = fmt.Errorf("opening template stop script (%q): %v", StopScriptFileName, err)
		return "", err
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		err = fmt.Errorf("reading template stop script (%q): %v", StopScriptFileName, err)
		return "", err
	}

	return b.String(), nil
}

func LoggingConfiguration(installPath string) (string, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "log4j2.xml", installPath); err != nil {
		err = fmt.Errorf("generating logging configuration file: %w", err)
		return "", err
	}

	return b.String(), nil
}
