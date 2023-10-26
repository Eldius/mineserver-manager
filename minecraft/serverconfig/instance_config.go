package serverconfig

import (
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"path/filepath"
)

type InstallOpts struct {
	Start       *RuntimeParams
	SrvProps    *ServerProperties
	Dest        string
	VersionName string
	VersionInfo *versions.VersionInfoResponse
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
			c.Start = GetDefaultRuntimeParams()
		}
		c.Start.Headless = headless
		return c
	}
}

func Headless() InstallOpt {
	return func(c *InstallOpts) *InstallOpts {
		if c.Start == nil {
			c.Start = GetDefaultRuntimeParams()
		}
		c.Start.Headless = true
		return c
	}
}

func NewInstallOpts(cfgs ...InstallOpt) *InstallOpts {
	cfg := &InstallOpts{
		Start:       GetDefaultRuntimeParams(),
		SrvProps:    utils.Must(GetDefaultServerProperties()),
		Dest:        "./minecraft",
		VersionName: "latest",
		VersionInfo: nil,
	}

	for _, c := range cfgs {
		cfg = c(cfg)
	}
	return cfg
}
