package minecraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/minecraft/config"
	"github.com/eldius/mineserver-manager/minecraft/installer"
	"github.com/eldius/mineserver-manager/minecraft/model"
	"github.com/eldius/mineserver-manager/minecraft/mojang"
	"github.com/eldius/mineserver-manager/minecraft/provisioner"
	"github.com/eldius/mineserver-manager/utils"
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

	Downloader     installer.Downloader
	RuntimeManager installer.RuntimeManager
	Provisioner    provisioner.Provisioner
	Flavor         installer.ServerFlavor
}

type InstallServiceOpt func(config *InstallServiceConfig)

type Installer interface {
	Install(ctx context.Context, configs ...config.InstanceOpt) error
}

type vanillaInstaller struct {
	cfg InstallServiceConfig
	d   installer.Downloader
	r   installer.RuntimeManager
	p   provisioner.Provisioner
	f   installer.ServerFlavor
}

// NewInstallService creates a new installer
func NewInstallService(configs ...InstallServiceOpt) Installer {
	cfg := &InstallServiceConfig{
		Timeout: 30 * time.Second,
	}
	for _, c := range configs {
		c(cfg)
	}

	if cfg.Downloader == nil {
		cfg.Downloader = installer.NewDownloader(cfg.DownloadTimeout)
	}
	if cfg.RuntimeManager == nil {
		cfg.RuntimeManager = installer.NewRuntimeManager(cfg.DownloadTimeout)
	}
	if cfg.Provisioner == nil {
		cfg.Provisioner = provisioner.NewProvisioner()
	}
	if cfg.Flavor == nil {
		cfg.Flavor = installer.NewVanillaFlavor(mojang.NewClient(mojang.WithTimeout(cfg.Timeout)))
	}

	return &vanillaInstaller{
		cfg: *cfg,
		d:   cfg.Downloader,
		r:   cfg.RuntimeManager,
		p:   cfg.Provisioner,
		f:   cfg.Flavor,
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

	info, err := i.f.GetVersionInfo(ctx, cfg.VersionName)
	if err != nil {
		return fmt.Errorf("getting version info for %s: %w", cfg.VersionName, err)
	}

	if err := i.p.CreateServerProperties(cfg.AbsoluteDestPath(), cfg.SrvProps); err != nil {
		return fmt.Errorf("creating server properties file: %w", err)
	}

	sf, err := i.d.DownloadServer(ctx, info.DownloadURL, info.SHA1, cfg.AbsoluteDestPath())
	if err != nil {
		return fmt.Errorf("downloading server file: %w", err)
	}

	log.With("server_file", sf).DebugContext(ctx, "Dowloaded server file")

	if _, err := i.r.InstallJava(ctx, filepath.Join(cfg.AbsoluteDestPath(), "java"), info.JavaVersion, runtime.GOARCH, runtime.GOOS); err != nil {
		return fmt.Errorf("installing jdk: %w", err)
	}

	if err := i.p.CreateStartScript(cfg.AbsoluteDestPath(),
		provisioner.WithHeadless(cfg.Headless),
		provisioner.WithJDKPath("${INSTALL_PATH}/java/jdk/bin"),
		provisioner.WithMemLimit(cfg.MemoryOpt),
		provisioner.WithServerFile("server.jar"),
		provisioner.WithLogConfigFile(cfg.AddLogConfig),
	); err != nil {
		return fmt.Errorf("creating start script: %w", err)
	}

	if err := i.p.CreateStopScript(cfg.AbsoluteDestPath()); err != nil {
		return fmt.Errorf("creating stop script: %w", err)
	}

	if cfg.AddLogConfig {
		if err := i.p.CreateLoggingConfig(cfg.AbsoluteDestPath(), cfg.AbsoluteDestPath()); err != nil {
			return fmt.Errorf("generating log config file: %w", err)
		}
	}

	if err := i.p.CreateEula(cfg.AbsoluteDestPath(), config.DefaultEulaValue); err != nil {
		return fmt.Errorf("creating eula.txt file: %w", err)
	}

	if err := i.createVersionFile(ctx, cfg.AbsoluteDestPath(), *cfg, info); err != nil {
		return fmt.Errorf("creating version file: %w", err)
	}

	// Whitelist requires mojang API specifically for UUID lookups
	// For now we only support it for vanilla if the flavor provides a client or we keep using mojang client directly.
	// Purpur might support it too if they use same UUIDs.
	if err := i.createWhitelistFile(ctx, *cfg); err != nil {
		return fmt.Errorf("creating whitelist file: %w", err)
	}

	return nil
}

func (i *vanillaInstaller) createWhitelistFile(_ context.Context, opts config.InstanceOpts) error {
	if !opts.HasWhitelist() {
		return nil
	}

	f, err := os.Create(filepath.Join(opts.Dest, "whitelist.json"))
	if err != nil {
		return fmt.Errorf("creating whitelist file: %w", err)
	}
	defer f.Close()

	// Direct use of mojang client for whitelist for now
	c := mojang.NewClient(mojang.WithTimeout(i.cfg.Timeout))
	usrs, err := c.GetUsersInfo(opts.WhitelistUsernames...)
	if err != nil {
		return fmt.Errorf("getting users info: %w", err)
	}

	if err := json.NewEncoder(f).Encode(usrs); err != nil {
		return fmt.Errorf("writing whitelist file: %w", err)
	}

	return nil
}

func (i *vanillaInstaller) createVersionFile(_ context.Context, destFolder string, opts config.InstanceOpts, info *installer.FlavorVersionInfo) error {

	f, err := os.Create(filepath.Join(destFolder, cfg.VersionsFileName))
	if err != nil {
		return fmt.Errorf("creating version.json file: %w", err)
	}
	defer f.Close()

	verInfo := cfg.GetVersionInfo()

	return json.NewEncoder(f).Encode(&model.VersionsInfo{
		JavaVersion: info.JavaVersion,
		MineVersion: info.Version,
		MineFlavour: i.f.Name(),
		CliVersion: model.CliVersion{
			Version:   verInfo.Version,
			Commit:    verInfo.Commit,
			BuildDate: verInfo.BuildDate,
		},
	})
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

func WithDownloader(d installer.Downloader) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Downloader = d
	}
}

func WithRuntimeManager(r installer.RuntimeManager) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.RuntimeManager = r
	}
}

func WithProvisioner(p provisioner.Provisioner) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Provisioner = p
	}
}

func WithFlavor(f installer.ServerFlavor) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Flavor = f
	}
}

func WithInstanceOpts(opts ...config.InstanceOpt) InstallServiceOpt {
	return func(cfg *InstallServiceConfig) {
		cfg.Instance = config.NewInstanceOpts(opts...)
	}
}
