package vanilla

import (
	"fmt"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/internal/vanilla/versions"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type InstallerConfig struct {
	Timeout time.Duration
}

type InstallerOpt func(config *InstallerConfig) *InstallerConfig

type Installer struct {
	cfg InstallerConfig
}

// NewInstaller creates a new client
func NewInstaller(configs ...InstallerOpt) *Installer {
	cfg := &InstallerConfig{
		Timeout: 1 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}
	return &Installer{
		cfg: *cfg,
	}
}

// InstallWithConfig installs selected version
func (i *Installer) InstallWithConfig(configs ...InstallCfg) error {
	cfg := installSetup(configs)

	c := versions.NewClient(versions.WithTimeout(i.cfg.Timeout))

	ver, err := c.ListVersions()
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		return err
	}

	v, err := ver.GetVersion(cfg.VersionName)
	if err != nil {
		err = fmt.Errorf("getting version from online versions list for name '%s': %w", cfg.VersionName, err)
		return err
	}

	cfg.v, err = c.GetVersionInfo(*v)
	if err != nil {
		err = fmt.Errorf("getting version info for name '%s': %w", cfg.VersionName, err)
		return err
	}

	sf, err := i.DownloadServer(*cfg.v, cfg.Dest)
	if err != nil {
		err = fmt.Errorf("getting version info to install: %w", err)
		return err
	}

	log.Printf("server file: %s", sf)
	return err
}

func installSetup(cfgs []InstallCfg) *InstallConfig {
	cfg := &InstallConfig{
		Start:       GetDefaultScriptParams(),
		SrvProps:    utils.Must(GetDefaultServerProperties()),
		Dest:        "./minecraft",
		VersionName: "latest",
		v:           nil,
	}

	for _, c := range cfgs {
		cfg = c(cfg)
	}
	return cfg
}

// InstallServer installs selected version
func (i *Installer) InstallServer(v versions.VersionInfoResponse, cfg *InstallConfig) error {
	sf, err := i.DownloadServer(v, cfg.Dest)
	if err != nil {
		err = fmt.Errorf("getting version info to install: %w", err)
		return err
	}

	log.Printf("server file: %s", sf)
	return err
}

// DownloadServer downloads server file
func (i *Installer) DownloadServer(v versions.VersionInfoResponse, dest string) (string, error) {
	client := i.httpInstaller()
	destFile := filepath.Join(dest, utils.GetFileName(v.Downloads.Server.URL))
	if err := utils.DownloadFile(&client, v.Downloads.Server.URL, destFile); err != nil {
		err = fmt.Errorf("getting version info: %w", err)
		return "", err
	}

	return destFile, nil
}

// StartScript generates the start script
func (i *Installer) StartScript(s StartupParams, dest string) error {
	destFile := filepath.Join(dest, "start.sh")

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server dest file: %w", err)
		return err
	}

	scp, err := s.ToScript()
	if err != nil {
		err = fmt.Errorf("generating start script content: %w", err)
		return err
	}

	if _, err := f.Write([]byte(scp)); err != nil {
		err = fmt.Errorf("writing start script to file: %w", err)
		return err
	}
	return nil
}

func (i *Installer) httpInstaller() http.Client {
	return http.Client{Timeout: i.cfg.Timeout}
}

func WithTimeout(t time.Duration) InstallerOpt {
	return func(cfg *InstallerConfig) *InstallerConfig {
		cfg.Timeout = t
		return cfg
	}
}
