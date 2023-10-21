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

type versionInfo struct {
	BuildDate  string
	Version    string
	CommitDate string
	Commit     string
}

func VersionInfo() {
	_ = template.Must(template.New("version").Parse(`---
version:     {{.Version}}
commit:      {{.Commit}}
commit date: {{.CommitDate}}
build date:  {{.BuildDate}}
`)).Execute(os.Stdout, versionInfo{
		BuildDate:  BuildDate,
		Version:    Version,
		CommitDate: CommitDate,
		Commit:     Commit,
	})
}
