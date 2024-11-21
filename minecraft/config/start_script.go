package config

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

var (
	//go:embed templates/**
	tmplFiles embed.FS
)

type StartupOptions struct {
	ServerFile    string
	JDKPath       string
	MemLimit      string
	LogConfigFile bool
	Headless      bool
}

type Option func(*StartupOptions)

func WithServerFile(serverFile string) Option {
	return func(o *StartupOptions) {
		if serverFile != "" {
			o.ServerFile = serverFile
		}
	}
}

func WithJDKPath(jdkPath string) Option {
	return func(o *StartupOptions) {
		if jdkPath != "" {
			o.JDKPath = jdkPath
		}
	}
}

func WithMemLimit(memLimit string) Option {
	return func(o *StartupOptions) {
		if memLimit != "" {
			o.MemLimit = memLimit
		}
	}
}

func WithHeadless(headless bool) Option {
	return func(o *StartupOptions) {
		o.Headless = headless
	}
}

func WithLogConfigFile(logConfigFile bool) Option {
	return func(o *StartupOptions) {
		o.LogConfigFile = logConfigFile
	}
}

func StartScript(opts ...Option) (string, error) {
	options := defaultStartupOptions()
	for _, o := range opts {
		o(options)
	}

	var b bytes.Buffer
	tmpl, err := template.ParseFS(tmplFiles, "templates/start.sh")
	if err != nil {
		return "", fmt.Errorf("parsing template file: %w", err)
	}
	if err := tmpl.Execute(&b, options); err != nil {
		return "", fmt.Errorf("generating start script: %w", err)
	}
	return b.String(), nil
}

func defaultStartupOptions() *StartupOptions {
	return &StartupOptions{
		ServerFile: "server.jar",
	}
}
