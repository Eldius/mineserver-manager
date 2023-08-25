package vanilla

import (
	"fmt"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/java"
	"github.com/eldius/mineserver-manager/vanilla/serverconfig"
	"github.com/eldius/mineserver-manager/vanilla/versions"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type InstallServiceConfig struct {
	Timeout         time.Duration
	DownloadTimeout time.Duration
	TargetFolder    string
}

type InstallServiceOpt func(config *InstallServiceConfig) *InstallServiceConfig

type InstallService struct {
	cfg InstallServiceConfig
}

// NewInstallService creates a new client
func NewInstallService(configs ...InstallServiceOpt) *InstallService {
	cfg := &InstallServiceConfig{
		Timeout: 30 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}
	return &InstallService{
		cfg: *cfg,
	}
}

// Install installs selected version
func (i *InstallService) Install(configs ...InstallCfg) error {
	cfg := installSetup(configs)

	log := logger.GetLogger().With("action", "install_server", "version_name", cfg.VersionName)

	c := versions.NewClient(versions.WithTimeout(i.cfg.Timeout))

	ver, err := c.ListVersions()
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		log.With("error", err).Error("Failed to list available versions")
		return err
	}

	v, err := ver.GetVersion(cfg.VersionName)
	if err != nil {
		err = fmt.Errorf("getting version from online versions list for name '%s': %w", cfg.VersionName, err)
		log.With("error", err).Error("Failed to get version for name %s", cfg.VersionName)
		return err
	}

	log = log.With("version", v.ID, "version_type", v.Type)

	cfg.v, err = c.GetVersionInfo(*v)
	if err != nil {
		err = fmt.Errorf("getting version info for name '%s': %w", cfg.VersionName, err)
		log.With("error", err).Error("Failed to fetch version info for '%s (%s)'", v.ID, cfg.VersionName)
		return err
	}

	sf, err := i.DownloadServer(*cfg.v, cfg.Dest)
	if err != nil {
		err = fmt.Errorf("downloading server file: %w", err)
		log.With("error", err).Error("Failed to download server file")
		return err
	}

	log.With("server_file", sf).Debug("Dowloaded server file")

	jdk, err := java.DownloadJDK(cfg.v.JavaVersion.MajorVersion, runtime.GOARCH, runtime.GOOS, i.cfg.DownloadTimeout)
	if err != nil {
		err = fmt.Errorf("downloading jdk package: %w", err)
		log.With("error", err).Error("Failed to download jdk")
		return err
	}

	if err = utils.UnpackTarGZ(jdk, cfg.Dest); err != nil {
		err = fmt.Errorf("unpacking jdk package: %w", err)
		log.Error("Failed to unpack JDK package: %v", err)
		return err
	}
	return err
}

func installSetup(cfgs []InstallCfg) *InstallConfig {
	cfg := &InstallConfig{
		Start:       serverconfig.GetDefaultScriptParams(),
		SrvProps:    utils.Must(serverconfig.GetDefaultServerProperties()),
		Dest:        "./minecraft",
		VersionName: "latest",
		v:           nil,
	}

	for _, c := range cfgs {
		cfg = c(cfg)
	}
	return cfg
}

// DownloadServer downloads server file
func (i *InstallService) DownloadServer(v versions.VersionInfoResponse, dest string) (string, error) {
	destFile := filepath.Join(dest, utils.GetFileName(v.Downloads.Server.URL))
	if err := utils.DownloadFile(i.cfg.DownloadTimeout, v.Downloads.Server.URL, destFile); err != nil {
		err = fmt.Errorf("getting version info: %w", err)
		return "", err
	}

	if err := utils.ValidateFileIntegrity(destFile, v.Downloads.Server.SHA1); err != nil {
		return "", err
	}

	return destFile, nil
}

// StartScript generates the start script
func (i *InstallService) StartScript(s serverconfig.StartupParams, dest string) error {
	destFile := filepath.Join(dest, "start.sh")

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
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

func (i *InstallService) httpInstaller() http.Client {
	return http.Client{Timeout: i.cfg.Timeout}
}

func WithTimeout(t time.Duration) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) *InstallServiceConfig {
		cfg.Timeout = t
		return cfg
	}
}

func WithDownloadTimeout(t time.Duration) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) *InstallServiceConfig {
		cfg.DownloadTimeout = t
		return cfg
	}
}
