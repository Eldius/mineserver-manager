package config

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
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

func readTemplateFile(filename string) (fs.File, error) {
	return templateFiles.Open(filename)
}

func StopScript() (string, error) {
	f, err := readTemplateFile(filepath.Join("templates", StopScriptFileName))
	if err != nil {
		err = fmt.Errorf("opening template stop script (%q): %v", StopScriptFileName, err)
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	fc, err := io.ReadAll(f)
	if err != nil {
		err = fmt.Errorf("reading stop script (%q): %v", StopScriptFileName, err)
		return "", err
	}

	return string(fc), nil
}

func LoggingConfiguration(logfileDestDir string) (string, error) {
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "log4j2.xml", logfileDestDir); err != nil {
		err = fmt.Errorf("generating logging configuration file: %w", err)
		return "", err
	}

	return b.String(), nil
}
