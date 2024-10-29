package model

import (
	"fmt"
	"github.com/asdine/storm/v3"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/google/uuid"
	"time"
)

var db *storm.DB

type RuntimeParams struct {
	Xmx           string
	Xms           string
	LogConfigFile bool
	Headless      bool
}

type Instance struct {
	ID               string `storm:"index"`
	Name             string `storm:"unique"`
	Path             string `storm:"unique"`
	InstallDate      time.Time
	RuntimeParams    serverconfig.RuntimeParams
	ServerProperties serverconfig.ServerProperties
}

func NewInstance(name, path string, runtimeParams serverconfig.RuntimeParams, serverProperties serverconfig.ServerProperties) *Instance {
	return &Instance{
		ID:               uuid.New().String(),
		Name:             name,
		Path:             path,
		InstallDate:      time.Now(),
		RuntimeParams:    runtimeParams,
		ServerProperties: serverProperties,
	}
}

func Persist(i *Instance) (*Instance, error) {
	if db == nil {
		openDB()
	}

	if err := db.Save(i); err != nil {
		err = fmt.Errorf("saving instance to db: %w", err)
		return nil, err
	}

	return i, nil
}

func openDB() {
	var err error
	db, err = storm.Open(config.GetAppHomePath())
	if err != nil {
		err = fmt.Errorf("opening db file: %w", err)
		panic(err)
	}
}

type CliVersion struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

type VersionsInfo struct {
	JavaVersion int        `json:"java_version"`
	MineVersion string     `json:"mine_version"`
	CliVersion  CliVersion `json:"cli_version"`
}

type WhitelistRecord struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
