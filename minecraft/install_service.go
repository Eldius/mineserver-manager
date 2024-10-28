package minecraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/java"
	"github.com/eldius/mineserver-manager/minecraft/model"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"github.com/eldius/properties"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	ErrChecksumValidationFailed = utils.ErrChecksumValidationFailed
)

type InstallServiceConfig struct {
	Timeout         time.Duration
	DownloadTimeout time.Duration
	TargetFolder    string
}

type InstallServiceOpt func(config *InstallServiceConfig) *InstallServiceConfig

type Installer interface {
	Install(ctx context.Context, configs ...serverconfig.InstallOpt) error
	DownloadServer(ctx context.Context, v versions.VersionInfoResponse, dest string) (string, error)
	CreateStartScript(s serverconfig.RuntimeParams, dest string) error
	CreateServerProperties(cfg *serverconfig.InstallOpts) error
	Eula(dest string) (string, error)
}

type vanillaInstaller struct {
	cfg InstallServiceConfig
}

// NewInstallService creates a new client
func NewInstallService(configs ...InstallServiceOpt) Installer {
	cfg := &InstallServiceConfig{
		Timeout: 30 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}
	return &vanillaInstaller{
		cfg: *cfg,
	}
}

// Install installs selected version
func (i *vanillaInstaller) Install(ctx context.Context, configs ...serverconfig.InstallOpt) error {
	cfg := serverconfig.NewInstallOpts(configs...)

	log := logger.GetLogger().With("action", "install_server", "version_name", cfg.VersionName)

	c := versions.NewClient(versions.WithTimeout(i.cfg.Timeout))

	if err := os.MkdirAll(cfg.AbsoluteDestPath(), os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			err = fmt.Errorf("creating destination folder: %w", err)
			log.With("error", err).ErrorContext(ctx, "Failed to create destination folder")
		}
		log.Debug("destination already exists")
	}

	fmt.Printf("#####################\nInstalling server\n----------------------\nversion: %s\nserver properties:\n%s\n#####################\n\n", cfg.VersionName, cfg.ServerPropertiesString())

	ver, err := c.ListVersions(ctx)
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to list available versions")
		return err
	}

	v, err := ver.GetVersion(cfg.VersionName)
	if err != nil {
		err = fmt.Errorf("getting version from online versions list for name '%s': %w", cfg.VersionName, err)
		log.With("error", err, "version", cfg.VersionName).ErrorContext(ctx, "Failed to get version for name")
		return err
	}

	log = log.With("version", v.ID, "version_type", v.Type)

	cfg.VersionInfo, err = c.GetVersionInfo(ctx, *v)
	if err != nil {
		err = fmt.Errorf("getting version info for name '%s': %w", cfg.VersionName, err)
		log.With("error", err).ErrorContext(ctx, "Failed to fetch version info for '%s (%s)'", v.ID, cfg.VersionName)
		return err
	}

	if err := i.CreateServerProperties(cfg); err != nil {
		err = fmt.Errorf("creating server properties file: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create server properties file")
		return err
	}

	sf, err := i.DownloadServer(ctx, *cfg.VersionInfo, cfg.AbsoluteDestPath())
	if err != nil {
		err = fmt.Errorf("downloading server file: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to download server file")
		return err
	}

	log.With("server_file", sf).DebugContext(ctx, "Dowloaded server file")

	if _, err := java.Install(ctx, filepath.Join(cfg.AbsoluteDestPath(), "java"), cfg.VersionInfo.JavaVersion.MajorVersion, runtime.GOARCH, runtime.GOOS, i.cfg.DownloadTimeout); err != nil {
		err = fmt.Errorf("downloading jdk package: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to download jdk")
		return err
	}

	if err := i.CreateStartScript(*cfg.Start, cfg.AbsoluteDestPath()); err != nil {
		err = fmt.Errorf("creating start script: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create start script")
		return err
	}

	if err := i.CreateStopScript(*cfg.Start, cfg.AbsoluteDestPath()); err != nil {
		err = fmt.Errorf("creating stop script: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create stop script")
		return err
	}

	if cfg.Start.LogConfigFile {
		if err := i.CreateLoggingConfig(*cfg.Start, cfg.AbsoluteDestPath()); err != nil {
			err = fmt.Errorf("generating log config file: %w", err)
			log.With("error", err).ErrorContext(ctx, "Failed to create log4j2.xml file")
			return err
		}
	}

	if _, err := i.Eula(cfg.AbsoluteDestPath()); err != nil {
		err = fmt.Errorf("creating eula.txt file: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create eula.txt file")
		return err
	}
	if err := i.createVersionFile(ctx, cfg.AbsoluteDestPath(), *cfg); err != nil {
		err = fmt.Errorf("creating version file: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create versions.json file")
		return err
	}
	return err
}

func (i *vanillaInstaller) createVersionFile(_ context.Context, destFolder string, opts serverconfig.InstallOpts) error {

	f, err := os.Create(filepath.Join(destFolder, config.VersionsFileName))
	if err != nil {
		err = fmt.Errorf("creating version.json file: %w", err)
		return err
	}

	info := config.GetVersionInfo()

	return json.NewEncoder(f).Encode(&model.VersionsInfo{
		JavaVersion: opts.VersionInfo.JavaVersion.MajorVersion,
		MineVersion: opts.VersionInfo.ID,
		CliVersion: model.CliVersion{
			Version:   info.Version,
			Commit:    info.Commit,
			BuildDate: info.BuildDate,
		},
	})
}

// DownloadServer downloads server file
func (i *vanillaInstaller) DownloadServer(ctx context.Context, v versions.VersionInfoResponse, dest string) (string, error) {
	destFile := filepath.Join(dest, utils.GetFileName(v.Downloads.Server.URL))
	if err := utils.DownloadFile(ctx, i.cfg.DownloadTimeout, v.Downloads.Server.URL, destFile); err != nil {
		err = fmt.Errorf("downloading server file: %w", err)
		return "", err
	}

	if err := utils.ValidateFileIntegrity(ctx, destFile, v.Downloads.Server.SHA1); err != nil {
		return "", err
	}

	return destFile, nil
}

func (i *vanillaInstaller) CreateServerProperties(cfg *serverconfig.InstallOpts) error {
	destFile := filepath.Join(cfg.AbsoluteDestPath(), "server.properties")
	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server properties file: %w", err)
		return err
	}

	if err := properties.NewEncoder(f).Encode(cfg.SrvProps); err != nil {
		err = fmt.Errorf("encoding server properties content to file: %w", err)
		return err
	}

	return nil
}

// CreateStartScript generates the start script
func (i *vanillaInstaller) CreateStartScript(s serverconfig.RuntimeParams, dest string) error {
	destFile := filepath.Join(dest, serverconfig.StartScriptFileName)

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
		return err
	}

	scp, err := s.RenderStartScript()
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

// CreateStopScript generates the stop server script
func (i *vanillaInstaller) CreateStopScript(s serverconfig.RuntimeParams, dest string) error {
	destFile := filepath.Join(dest, serverconfig.StopScriptFileName)

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server stop script: %w", err)
		return err
	}

	scp, err := s.RenderStopScript()
	if err != nil {
		err = fmt.Errorf("generating stop script content: %w", err)
		return err
	}

	if _, err := f.Write([]byte(scp)); err != nil {
		err = fmt.Errorf("writing stop script to file: %w", err)
		return err
	}
	return nil
}

// CreateLoggingConfig generates log4j2.xml logging configuration file
func (i *vanillaInstaller) CreateLoggingConfig(s serverconfig.RuntimeParams, dest string) error {
	f, err := os.OpenFile(filepath.Join(dest, "log4j2.xml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
		return err
	}

	logf, err := s.LoggingConfiguration(dest)
	if err != nil {
		err = fmt.Errorf("generating logging configuration file content: %w", err)
		return err
	}

	if _, err := f.Write([]byte(logf)); err != nil {
		err = fmt.Errorf("writing logging configuration file to file: %w", err)
		return err
	}
	return nil
}

func (i *vanillaInstaller) Eula(dest string) (string, error) {
	eulaPath := filepath.Join(dest, "eula.txt")
	f, err := os.Create(eulaPath)
	if err != nil {
		err = fmt.Errorf("creating eula file: %w", err)
		return "", err
	}
	if err := properties.NewEncoder(f).Encode(serverconfig.DefaultEulaValue); err != nil {
		err = fmt.Errorf("writing eula contents to file: %w", err)
		return "", err
	}

	return eulaPath, nil
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
