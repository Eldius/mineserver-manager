package minecraft

import (
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"path/filepath"
)

type InstallOpts struct {
	Start       *serverconfig.StartupParams
	SrvProps    *serverconfig.ServerProperties
	Dest        string
	VersionName string
	v           *versions.VersionInfoResponse
}

func (o InstallOpts) AbsoluteDestPath() string {
	d, err := filepath.Abs(o.Dest)
	if err != nil {
		return o.Dest
	}
	return d
}

type InstallOpt func(*InstallOpts) *InstallOpts

func WithVersion(v string) InstallOpt {
	return func(c *InstallOpts) *InstallOpts {
		c.VersionName = v
		return c
	}
}

func ToDestinationFolder(t string) InstallOpt {
	return func(c *InstallOpts) *InstallOpts {
		c.Dest = t
		return c
	}
}

func WithHeadlessConfig(headless bool) InstallOpt {
	return func(c *InstallOpts) *InstallOpts {
		if c.Start == nil {
			c.Start = serverconfig.GetDefaultScriptParams()
		}
		c.Start.Headless = headless
		return c
	}
}

func Headless() InstallOpt {
	return func(c *InstallOpts) *InstallOpts {
		if c.Start == nil {
			c.Start = serverconfig.GetDefaultScriptParams()
		}
		c.Start.Headless = true
		return c
	}
}
