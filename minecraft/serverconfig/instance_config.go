package serverconfig

import (
	"bytes"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"gopkg.in/yaml.v3"
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

func (o InstallOpts) ServerPropertiesString() string {
	var buffer bytes.Buffer
	_ = yaml.NewEncoder(&buffer).Encode(o.SrvProps)
	return buffer.String()
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

func WithServerPropsMotd(m string) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.Motd = m
		return s
	}
}

func WithServerPropsLevelName(n string) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.LevelName = n
		return s
	}
}

func WithServerPropsServerPort(p int) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.ServerPort = p
		return s
	}
}

// WithServerPropsRconEnabled enables RCON protocol configuration
// 'port' to be used for this protocol
// 'pass' define the RCON password
func WithServerPropsRconEnabled(port int, pass string) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.RconPort = port
		s.SrvProps.EnableRcon = true
		s.SrvProps.RconPassword = pass
		return s
	}
}

// WithServerPropsRcon defines RCON protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
// 'pass' define the RCON password
func WithServerPropsRcon(port int, enabled bool, pass string) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.RconPort = port
		s.SrvProps.EnableRcon = enabled
		s.SrvProps.RconPassword = pass
		return s
	}
}

// WithServerPropsQuery defines Query protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
func WithServerPropsQuery(port int, enabled bool) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.QueryPort = port
		s.SrvProps.EnableQuery = enabled
		return s
	}
}

// WithServerPropsSeed defines level seed
func WithServerPropsSeed(seed string) InstallOpt {
	return func(s *InstallOpts) *InstallOpts {
		s.SrvProps.LevelSeed = seed
		return s
	}
}
