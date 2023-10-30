package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/minecraft"
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"time"
)

// installCmd installs a minecraft server instance
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a Minecraft server instance",
	Long:  `Installs a Minecraft server instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.GetLogger().With("headless", installOpts.Headless).Info("debugging")
		ctx := context.Background()
		client := minecraft.NewInstallService(minecraft.WithTimeout(
			config.GetMinecraftApiTimeout()),
			minecraft.WithDownloadTimeout(config.GetMinecraftDownloadTimeout()),
		)

		if installOpts.JustListVersions {
			if err := minecraft.ListVersions(ctx); err != nil {
				err = fmt.Errorf("listing available versions: %w", err)
				slog.ErrorContext(ctx, "failed to list available versions: %v", err)
				panic(err)
			}
		} else {
			opts := append(
				installOpts.ToServerPropertiesOpts(),
				serverconfig.WithVersion(installOpts.ServerVersion),
				serverconfig.ToDestinationFolder(installOpts.DestinationFolder),
				serverconfig.WithHeadlessConfig(installOpts.Headless),
			)
			if err := client.Install(ctx, opts...); err != nil {
				err = fmt.Errorf("installing server: %w", err)
				slog.ErrorContext(ctx, "failed to install server: %v", err)
				panic(err)
			}
		}
	},
}

type installCmdOpts struct {
	ServerVersion     string
	DestinationFolder string
	Headless          bool
	JustListVersions  bool

	Motd       string
	LevelName  string
	Seed       string
	ServerPort int
	QueryPort  int

	RconPort    int
	RconPass    string
	RconEnabled bool
}

func (o installCmdOpts) ToServerPropertiesOpts() []serverconfig.InstallOpt {
	var opts []serverconfig.InstallOpt
	if o.Motd != "" {
		opts = append(opts, serverconfig.WithServerPropsMotd(o.Motd))
	}
	if o.LevelName != "" {
		opts = append(opts, serverconfig.WithServerPropsLevelName(o.LevelName))
	}
	if o.Seed != "" {
		opts = append(opts, serverconfig.WithServerPropsSeed(o.Seed))
	}
	if o.RconEnabled {
		passwd, err := utils.PasswordPrompt()
		if err != nil {
			err = fmt.Errorf("rcon pass promtp: %w", err)
			panic(err)
		}
		if passwd == "" {
			err = errors.New("RCON password mustn't be empty when RCON is enabled")
			panic(err)
		}
		opts = append(opts, serverconfig.WithServerPropsRconEnabled(o.RconPort, passwd))
	}
	return opts
}

var (
	installOpts = installCmdOpts{}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installOpts.ServerVersion, "version", "latest", "Java Edition server version to be installed, ('latest' will install latest stable version)")
	installCmd.Flags().StringVar(&installOpts.DestinationFolder, "dest", ".", "Installation root directory (defaults to current directory)")
	installCmd.Flags().BoolVar(&installOpts.Headless, "headless", false, "Installation root directory (defaults to false)")
	installCmd.Flags().BoolVar(&installOpts.JustListVersions, "list", false, "Lists available versions to install")

	installCmd.Flags().StringVar(&installOpts.Motd, "motd", "A Minecraft Server", "Server name (defaults to 'A Minecraft Server')")
	installCmd.Flags().StringVar(&installOpts.LevelName, "level-name", "", "Level/map name")
	installCmd.Flags().StringVar(&installOpts.Seed, "seed", "", "Seed to be used to generate game map")

	installCmd.Flags().IntVar(&installOpts.ServerPort, "server-port", 25565, "Server port (defaults to 25565)")
	installCmd.Flags().IntVar(&installOpts.QueryPort, "query-port", 25566, "Server port (defaults to 25565)")

	installCmd.Flags().IntVar(&installOpts.RconPort, "rcon-port", 25575, "RCON server port (defaults to 25565)")
	installCmd.Flags().BoolVar(&installOpts.RconEnabled, "enable-rcon", false, "Enable RCON protocol")

	installCmd.Flags().Duration("download-timeout", 300*time.Second, "Download timeout configuration (defaults to 300s/5m)")
	if err := viper.BindPFlag(config.AppInstallDownloadTimeoutPropKey, installCmd.Flags().Lookup("download-timeout")); err != nil {
		err = fmt.Errorf("binding artifact download timeout property: %w", err)
		panic(err)
	}
}
