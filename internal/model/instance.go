package model

import (
	"github.com/google/uuid"
	"time"
)

type Instance struct {
	ID   string `storm:"index"`
	Name string `storm:"unique"`
	Path string `storm:"unique"`

	InstallDate      time.Time
	ServerProperties ServerProperties
}

type RemoteInstance struct {
	Instance
	Host string `storm:"index"`
	Port int64
	User string `storm:"index"`
}

func NewInstance(name, path string, serverProperties ServerProperties) *Instance {
	return &Instance{
		ID:               uuid.New().String(),
		Name:             name,
		Path:             path,
		InstallDate:      time.Now(),
		ServerProperties: serverProperties,
	}
}

type CliVersion struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

type MineFlavour string

const (
	MineFlavourVanilla MineFlavour = "vanilla"
	MineFlavourPurpur  MineFlavour = "purpur"
)

type VersionsInfo struct {
	CliVersion  CliVersion  `json:"cli_version"`
	MineFlavour MineFlavour `json:"mine_flavour"`
	MineVersion string      `json:"mine_version"`
	JavaVersion int         `json:"java_version"`
}

type WhitelistRecord struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
