package config

import (
	"bytes"
	"github.com/eldius/mineserver-manager/internal/utils"
	"github.com/eldius/mineserver-manager/minecraft/mojang"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

type InstanceOpts struct {
	SrvProps           *ServerProperties
	VersionInfo        *mojang.VersionInfoResponse
	Dest               string
	VersionName        string
	WhitelistUsernames []string
	MemoryOpt          string
	AddLogConfig       bool
	Headless           bool
}

func (o InstanceOpts) HasWhitelist() bool {
	return len(o.WhitelistUsernames) > 0
}

func (o InstanceOpts) AbsoluteDestPath() string {
	d, err := filepath.Abs(o.Dest)
	if err != nil {
		return o.Dest
	}
	return d
}

func (o InstanceOpts) ServerPropertiesString() string {
	var buffer bytes.Buffer
	_ = yaml.NewEncoder(&buffer).Encode(o.SrvProps)
	return buffer.String()
}

type InstanceOpt func(*InstanceOpts)

func WithVersion(v string) InstanceOpt {
	return func(c *InstanceOpts) {
		c.VersionName = v
	}
}

func WithWhitelistedUsers(users []string) InstanceOpt {
	return func(c *InstanceOpts) {
		if len(users) == 0 {
			return
		}
		c.WhitelistUsernames = users
	}
}

func ToDestinationFolder(t string) InstanceOpt {
	return func(c *InstanceOpts) {
		c.Dest = t
	}
}

func WithHeadlessConfig(headless bool) InstanceOpt {
	return func(c *InstanceOpts) {
		c.Headless = headless
	}
}

func Headless() InstanceOpt {
	return func(c *InstanceOpts) {
		c.Headless = true
	}
}

func WithMemoryLimit(memory string) InstanceOpt {
	return func(c *InstanceOpts) {
		c.MemoryOpt = memory
	}
}

func NewInstanceOpts(cfgs ...InstanceOpt) *InstanceOpts {
	cfg := &InstanceOpts{
		SrvProps:           utils.Must(DefaultServerProperties()),
		Dest:               "./minecraft",
		VersionName:        "latest",
		MemoryOpt:          "1g",
		VersionInfo:        nil,
		WhitelistUsernames: nil,
	}

	for _, c := range cfgs {
		c(cfg)
	}
	return cfg
}

func WithServerPropsMotd(m string) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.Motd = m
	}
}

func WithServerPropsLevelName(n string) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.LevelName = n
	}
}

func WithServerPropsServerPort(p int) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.ServerPort = p
	}
}

// WithServerPropsRconEnabled enables RCON protocol configuration
// 'port' to be used for this protocol
// 'pass' define the RCON password
func WithServerPropsRconEnabled(port int, pass string) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.RconPort = port
		s.SrvProps.EnableRcon = true
		s.SrvProps.RconPassword = pass
	}
}

// WithServerPropsRcon defines RCON protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
// 'pass' define the RCON password
func WithServerPropsRcon(port int, enabled bool, pass string) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.RconPort = port
		s.SrvProps.EnableRcon = enabled
		s.SrvProps.RconPassword = pass
	}
}

// WithServerPropsQuery defines Query protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
func WithServerPropsQuery(port int, enabled bool) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.QueryPort = port
		s.SrvProps.EnableQuery = enabled
	}
}

// WithServerPropsSeed defines level seed
func WithServerPropsSeed(seed string) InstanceOpt {
	return func(s *InstanceOpts) {
		s.SrvProps.LevelSeed = seed
	}
}
