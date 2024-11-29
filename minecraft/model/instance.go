package model

import (
	"fmt"
	"github.com/asdine/storm/v3"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/minecraft/config"
	"github.com/google/uuid"
	"time"
)

var db *storm.DB

type Instance struct {
	ID               string `storm:"index"`
	Name             string `storm:"unique"`
	Path             string `storm:"unique"`
	InstallDate      time.Time
	ServerProperties config.ServerProperties
}

func NewInstance(name, path string, serverProperties config.ServerProperties) *Instance {
	return &Instance{
		ID:               uuid.New().String(),
		Name:             name,
		Path:             path,
		InstallDate:      time.Now(),
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
	db, err = storm.Open(cfg.GetAppHomePath())
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

type MineFlavour string

const (
	MineFlavourVanilla MineFlavour = "vanilla"
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
