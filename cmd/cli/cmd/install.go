package cmd

import (
	"context"
	cfg "github.com/eldius/mineserver-manager/internal/config"
	"github.com/eldius/mineserver-manager/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

// installCmd installs a minecraft server instance
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a Minecraft server instance",
	Long:  `Installs a Minecraft server instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.GetLogger().With("headless", installOpts.Headless).Info("debugging")
		return runInstall(context.Background(), installOpts)
	},
}

type installCmdOpts struct {
	Flavor            string
	ServerVersion     string
	DestinationFolder string
	Headless          bool
	JustListVersions  bool

	Motd       string
	LevelName  string
	Seed       string
	ServerPort int
	QueryPort  int

	MemoryLimit string

	RconPort    int
	RconPass    string
	RconEnabled bool

	users []string
}

var (
	installOpts = installCmdOpts{}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installOpts.Flavor, "flavor", "vanilla", "Minecraft server flavor (vanilla, purpur)")
	installCmd.Flags().StringVar(&installOpts.ServerVersion, "version", "latest", "Java Edition server version to be installed, ('latest' will install latest stable version)")
	installCmd.Flags().StringVar(&installOpts.DestinationFolder, "dest", ".", "Installation root directory (defaults to current directory)")
	installCmd.Flags().BoolVar(&installOpts.Headless, "headless", false, "Installation root directory (defaults to false)")
	installCmd.Flags().BoolVar(&installOpts.JustListVersions, "list", false, "Lists available versions to install")
	installCmd.Flags().StringVar(&installOpts.MemoryLimit, "memory-limit", "1g", "Server memory limit")

	installCmd.Flags().StringVar(&installOpts.Motd, "motd", "A Minecraft Server", "Server name (defaults to 'A Minecraft Server')")
	installCmd.Flags().StringVar(&installOpts.LevelName, "level-name", "", "Level/map name")
	installCmd.Flags().StringVar(&installOpts.Seed, "seed", "", "Seed to be used to generate game map")

	installCmd.Flags().IntVar(&installOpts.ServerPort, "server-port", 25565, "Server port (defaults to 25565)")
	installCmd.Flags().IntVar(&installOpts.QueryPort, "query-port", 25566, "Server port (defaults to 25565)")

	installCmd.Flags().BoolVar(&installOpts.RconEnabled, "rcon-enabled", false, "Enable RCON protocol")
	installCmd.Flags().IntVar(&installOpts.RconPort, "rcon-port", 25575, "RCON server port (defaults to 25565)")
	installCmd.Flags().StringVar(&installOpts.RconPass, "rcon-passwd", "", "RCON password (it will be asked if empty)")

	installCmd.Flags().StringSliceVar(&installOpts.users, "whitelist-user", []string{}, "List of users to whitelist (optional)")

	installCmd.Flags().Duration("download-timeout", 300*time.Second, "Download timeout configuration (defaults to 300s/5m)")
	if err := viper.BindPFlag(cfg.AppInstallDownloadTimeoutPropKey, installCmd.Flags().Lookup("download-timeout")); err != nil {
		panic(err)
	}
}
