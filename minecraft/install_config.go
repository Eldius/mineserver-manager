package minecraft

import (
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/eldius/mineserver-manager/minecraft/versions"
)

type InstallConfig struct {
	Start       *serverconfig.StartupParams
	SrvProps    *serverconfig.ServerProperties
	Dest        string
	VersionName string
	v           *versions.VersionInfoResponse
}

type InstallCfg func(*InstallConfig) *InstallConfig

func WithVersion(v string) InstallCfg {
	return func(c *InstallConfig) *InstallConfig {
		c.VersionName = v
		return c
	}
}

func ToDestinationFolder(t string) InstallCfg {
	return func(c *InstallConfig) *InstallConfig {
		c.Dest = t
		return c
	}
}
