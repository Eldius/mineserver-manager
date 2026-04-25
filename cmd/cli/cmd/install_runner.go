package cmd

import (
	"context"
	"errors"
	"fmt"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/installer"
	"github.com/eldius/mineserver-manager/internal/minecraft"
	"github.com/eldius/mineserver-manager/internal/minecraft/config"
	"github.com/eldius/mineserver-manager/internal/mojang"
	"github.com/eldius/mineserver-manager/internal/utils"
	"log/slog"
)

func runInstall(ctx context.Context, opts installCmdOpts) error {
	var flavor installer.ServerFlavor
	switch opts.Flavor {
	case "vanilla":
		flavor = installer.NewVanillaFlavor(mojang.NewClient(mojang.WithTimeout(cfg.GetMinecraftApiTimeout())))
	case "purpur":
		return errors.New("purpur flavor not yet implemented")
	default:
		return fmt.Errorf("invalid flavor: %s", opts.Flavor)
	}

	if opts.JustListVersions {
		versions, err := flavor.ListVersions(ctx)
		if err != nil {
			return fmt.Errorf("listing available versions: %w", err)
		}
		for _, v := range versions {
			fmt.Printf("- %s\n", v)
		}
		return nil
	}

	client := minecraft.NewInstallService(
		minecraft.WithTimeout(cfg.GetMinecraftApiTimeout()),
		minecraft.WithDownloadTimeout(cfg.GetMinecraftDownloadTimeout()),
		minecraft.WithFlavor(flavor),
	)

	instanceOpts := append(
		opts.ToInstanceOpts(),
		config.WithVersion(opts.ServerVersion),
		config.ToDestinationFolder(opts.DestinationFolder),
		config.WithHeadlessConfig(opts.Headless),
		config.WithServerFlavour(opts.Flavor),
	)

	if err := client.Install(ctx, instanceOpts...); err != nil {
		return fmt.Errorf("installing server: %w", err)
	}

	return nil
}

func (o installCmdOpts) ToInstanceOpts() []config.InstanceOpt {
	opts := []config.InstanceOpt{config.WithMemoryLimit(o.MemoryLimit)}
	if o.Motd != "" {
		opts = append(opts, config.WithServerPropsMotd(o.Motd))
	}
	if o.LevelName != "" {
		opts = append(opts, config.WithServerPropsLevelName(o.LevelName))
	}
	if o.Seed != "" {
		opts = append(opts, config.WithServerPropsSeed(o.Seed))
	}
	if o.RconEnabled {
		if o.RconPass == "" {
			passwd, err := utils.PasswordPrompt()
			if err != nil {
				slog.With("error", err).Error("failed to prompt for password")
				panic(err)
			}
			o.RconPass = passwd
		}
		if o.RconPass == "" {
			panic(errors.New("RCON password mustn't be empty when RCON is enabled"))
		}
		opts = append(opts, config.WithServerPropsRconEnabled(o.RconPort, o.RconPass))
	}

	if len(o.users) > 0 {
		opts = append(opts, config.WithWhitelistedUsers(o.users))
	}

	return opts
}
