package provisioner

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/eldius/properties"
	"io"
	"io/fs"
	"os"
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

type Provisioner interface {
	CreateServerProperties(dest string, props *model.ServerProperties) error
	CreateStartScript(dest string, opts ...StartupOption) error
	CreateStopScript(dest string) error
	CreateLoggingConfig(dest string, logfileDestDir string) error
	CreateEula(dest string, eula *model.Eula) error
}

type vanillaProvisioner struct{}

func NewProvisioner() Provisioner {
	return &vanillaProvisioner{}
}

func (p *vanillaProvisioner) CreateServerProperties(dest string, props *model.ServerProperties) error {
	destFile := filepath.Join(dest, "server.properties")
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating server properties file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := properties.NewEncoder(f).Encode(props); err != nil {
		return fmt.Errorf("encoding server properties: %w", err)
	}

	return nil
}

func (p *vanillaProvisioner) CreateStartScript(dest string, opts ...StartupOption) error {
	script, err := StartScript(opts...)
	if err != nil {
		return err
	}
	return p.writeFile(filepath.Join(dest, StartScriptFileName), script, 0755)
}

func (p *vanillaProvisioner) CreateStopScript(dest string) error {
	script, err := StopScript()
	if err != nil {
		return err
	}
	return p.writeFile(filepath.Join(dest, StopScriptFileName), script, 0755)
}

func (p *vanillaProvisioner) CreateLoggingConfig(dest string, logfileDestDir string) error {
	config, err := LoggingConfiguration(logfileDestDir)
	if err != nil {
		return err
	}
	return p.writeFile(filepath.Join(dest, "log4j2.xml"), config, 0644)
}

func (p *vanillaProvisioner) CreateEula(dest string, eula *model.Eula) error {
	destFile := filepath.Join(dest, "eula.txt")
	f, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("creating eula file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := properties.NewEncoder(f).Encode(eula); err != nil {
		return fmt.Errorf("writing eula contents: %w", err)
	}
	return nil
}

func (p *vanillaProvisioner) writeFile(path string, content string, perm os.FileMode) error {
	return os.WriteFile(path, []byte(content), perm)
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
