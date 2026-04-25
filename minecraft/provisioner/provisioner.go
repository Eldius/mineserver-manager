package provisioner

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

func readTemplateFile(filename string) (fs.File, error) {
	return templateFiles.Open(filename)
}

type StartupOptions struct {
	ServerFile    string
	JDKPath       string
	MemLimit      string
	LogConfigFile bool
	Headless      bool
}

type StartupOption func(*StartupOptions)

func WithServerFile(serverFile string) StartupOption {
	return func(o *StartupOptions) {
		if serverFile != "" {
			o.ServerFile = serverFile
		}
	}
}

func WithJDKPath(jdkPath string) StartupOption {
	return func(o *StartupOptions) {
		if jdkPath != "" {
			o.JDKPath = jdkPath
		}
	}
}

func WithMemLimit(memLimit string) StartupOption {
	return func(o *StartupOptions) {
		if memLimit != "" {
			o.MemLimit = memLimit
		}
	}
}

func WithHeadless(headless bool) StartupOption {
	return func(o *StartupOptions) {
		o.Headless = headless
	}
}

func WithLogConfigFile(logConfigFile bool) StartupOption {
	return func(o *StartupOptions) {
		o.LogConfigFile = logConfigFile
	}
}

func StartScript(opts ...StartupOption) (string, error) {
	options := defaultStartupOptions()
	for _, o := range opts {
		o(options)
	}

	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "start.sh", options); err != nil {
		return "", fmt.Errorf("generating start script: %w", err)
	}
	return b.String(), nil
}

func defaultStartupOptions() *StartupOptions {
	return &StartupOptions{
		ServerFile: "server.jar",
	}
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
