package config

import (
	"os"
	"text/template"
)

var (
	BuildDate  string
	Version    string
	CommitDate string
	Commit     string
)

type VersionInfo struct {
	BuildDate  string
	Version    string
	CommitDate string
	Commit     string
}

func GetVersionInfo() VersionInfo {
	return VersionInfo{
		BuildDate:  BuildDate,
		Version:    Version,
		CommitDate: CommitDate,
		Commit:     Commit,
	}
}

func DisplayVersionInfo() {
	_ = template.Must(template.New("version").Parse(`---
version:     {{.Version}}
commit:      {{.Commit}}
commit date: {{.CommitDate}}
build date:  {{.BuildDate}}
`)).Execute(os.Stdout, GetVersionInfo())
}
