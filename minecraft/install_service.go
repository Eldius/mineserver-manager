package minecraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/java"
	"github.com/eldius/mineserver-manager/minecraft/config"
	"github.com/eldius/mineserver-manager/minecraft/model"
	"github.com/eldius/mineserver-manager/minecraft/mojang"
	"github.com/eldius/mineserver-manager/utils"
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
	Instance        *config.InstanceOpts
	Timeout         time.Duration
	DownloadTimeout time.Duration
	TargetFolder    string
}

type InstallServiceOpt func(config *InstallServiceConfig)

type Installer interface {
	Install(ctx context.Context, configs ...config.InstanceOpt) error
	DownloadServer(ctx context.Context, v mojang.VersionInfoResponse, dest string) (string, error)
	CreateStartScript(cfg *config.InstanceOpts) error
	CreateServerProperties(cfg *config.InstanceOpts) error
	Eula(dest string) (string, error)
}

type vanillaInstaller struct {
	cfg InstallServiceConfig
	c   mojang.Client
}

// NewInstallService creates a new installer
func NewInstallService(configs ...InstallServiceOpt) Installer {
	cfg := &InstallServiceConfig{
		Timeout: 30 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}

	return &vanillaInstaller{
		cfg: *cfg,
		c:   mojang.NewClient(mojang.WithTimeout(cfg.Timeout)),
	}
}

// Install installs selected version
func (i *vanillaInstaller) Install(ctx context.Context, configs ...config.InstanceOpt) error {
	cfg := config.NewInstanceOpts(configs...)

	log := logger.GetLogger().With("action", "install_server", "version_name", cfg.VersionName)

	if err := os.MkdirAll(cfg.AbsoluteDestPath(), os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			err = fmt.Errorf("creating destination folder: %w", err)
			log.With("error", err).ErrorContext(ctx, "Failed to create destination folder")
		}
		log.Debug("destination already exists")
	}

	fmt.Printf("#####################\nInstalling server\n----------------------\nversion: %s\nserver properties:\n%s\n#####################\n\n", cfg.VersionName, cfg.ServerPropertiesString())

	ver, err := i.c.ListVersions(ctx)
	if err != nil {
		err = fmt.Errorf("getting available versions: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to list available versions")
		return err
	}

	v, err := ver.GetVersion(cfg.VersionName)
	if err != nil {
		err = fmt.Errorf("getting online versions list for name '%s': %w", cfg.VersionName, err)
		log.With("error", err, "version", cfg.VersionName).ErrorContext(ctx, "Failed to get version for name")
		return err
	}

	log = log.With("version", v.ID, "version_type", v.Type)

	cfg.VersionInfo, err = i.c.GetVersionInfo(ctx, *v)
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

	if err := i.CreateStartScript(cfg); err != nil {
		err = fmt.Errorf("creating start script: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create start script")
		return err
	}

	if err := i.CreateStopScript(cfg.AbsoluteDestPath()); err != nil {
		err = fmt.Errorf("creating stop script: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create stop script")
		return err
	}

	if cfg.AddLogConfig {
		if err := i.CreateLoggingConfig(cfg.AbsoluteDestPath()); err != nil {
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

	if err := i.createWhitelistFile(ctx, *cfg); err != nil {
		err = fmt.Errorf("creating whitelist file: %w", err)
		log.With("error", err).ErrorContext(ctx, "Failed to create whitelist.json file")
		return err
	}

	return err
}

func (i *vanillaInstaller) createWhitelistFile(_ context.Context, opts config.InstanceOpts) error {
	if !opts.HasWhitelist() {
		return nil
	}

	f, err := os.Create(filepath.Join(opts.Dest, "whitelist.json"))
	if err != nil {
		err = fmt.Errorf("creating whitelist file: %w", err)
		return err
	}

	usrs, err := i.c.GetUsersInfo(opts.WhitelistUsernames...)
	if err != nil {
		err = fmt.Errorf("getting users info: %w", err)
		return err
	}

	if err := json.NewEncoder(f).Encode(usrs); err != nil {
		err = fmt.Errorf("writing whitelist file: %w", err)
		return err
	}

	return nil
}

func (i *vanillaInstaller) createVersionFile(_ context.Context, destFolder string, opts config.InstanceOpts) error {

	f, err := os.Create(filepath.Join(destFolder, cfg.VersionsFileName))
	if err != nil {
		err = fmt.Errorf("creating version.json file: %w", err)
		return err
	}

	info := cfg.GetVersionInfo()

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
func (i *vanillaInstaller) DownloadServer(ctx context.Context, v mojang.VersionInfoResponse, dest string) (string, error) {
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

func (i *vanillaInstaller) CreateServerProperties(cfg *config.InstanceOpts) error {
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
func (i *vanillaInstaller) CreateStartScript(cfg *config.InstanceOpts) error {

	destFile := filepath.Join(cfg.AbsoluteDestPath(), config.StartScriptFileName)

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
		return err
	}

	script, err := config.StartScript(
		config.WithHeadless(cfg.Headless),
		config.WithJDKPath("${INSTALL_PATH}/java/jdk/bin"),
		config.WithMemLimit(cfg.MemoryOpt),
		config.WithServerFile("server.jar"),
		config.WithLogConfigFile(cfg.AddLogConfig),
	)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
		return err
	}

	if _, err := f.Write([]byte(script)); err != nil {
		err = fmt.Errorf("writing start script to file: %w", err)
		return err
	}
	return nil
}

// CreateStopScript generates the stop server script
func (i *vanillaInstaller) CreateStopScript(dest string) error {
	destFile := filepath.Join(dest, config.StopScriptFileName)

	f, err := os.OpenFile(destFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server stop script: %w", err)
		return err
	}

	scp, err := config.StopScript()
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
func (i *vanillaInstaller) CreateLoggingConfig(dest string) error {
	f, err := os.OpenFile(filepath.Join(dest, "log4j2.xml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("creating server startup script: %w", err)
		return err
	}

	logf, err := config.LoggingConfiguration(dest)
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
	if err := properties.NewEncoder(f).Encode(config.DefaultEulaValue); err != nil {
		err = fmt.Errorf("writing eula contents to file: %w", err)
		return "", err
	}

	return eulaPath, nil
}

func WithTimeout(t time.Duration) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Timeout = t
	}
}

func WithDownloadTimeout(t time.Duration) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.DownloadTimeout = t
	}
}

func WithInstanceOpts(opts ...config.InstanceOpt) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Instance = config.NewInstanceOpts(opts...)
	}
}
