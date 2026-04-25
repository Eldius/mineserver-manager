package utils

import (
	"context"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/netip"
	"testing"
	"time"
)

const (
	openSshServerImage = "lscr.io/linuxserver/openssh-server:latest"
	exposedPort        = "2222/tcp"
)

func TestSSHConnection(t *testing.T) {

	sshHost, cancel := startSshServer(t)
	defer cancel()

	ip, err := sshHost.ContainerIP(t.Context())
	assert.NoError(t, err)

	t.Log("ssh server ip:", ip)

	ctx, cancel1 := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel1()

	assert.NoError(t, Execute(ctx, "127.0.0.1", "2222"))
}

type myLogConsumer struct {
	testcontainers.LogConsumer
	t *testing.T
}

func (c *myLogConsumer) Accept(log testcontainers.Log) {
	c.t.Logf("[container] %s", log.Content)
}

func startSshServer(t *testing.T) (testcontainers.Container, context.CancelFunc) {
	t.Helper()

	ctx := t.Context()

	req := containerRequest(t)
	sshHost, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	assert.NoError(t, err)

	ip, err := sshHost.ContainerIP(ctx)
	assert.NoError(t, err)
	t.Logf("SSH Host: %s", ip)

	return sshHost, func() {
		testcontainers.CleanupContainer(t, sshHost)
	}
}

func containerRequest(t *testing.T) testcontainers.ContainerRequest {
	t.Helper()

	log := &myLogConsumer{
		t: t,
	}
	return testcontainers.ContainerRequest{
		Image: openSshServerImage,
		Env: map[string]string{
			"PUID":            "1000",
			"PGID":            "1000",
			"TZ":              "GMT",
			"SUDO_ACCESS":     "true",
			"PASSWORD_ACCESS": "true",
			"USER_PASSWORD":   sshPass,
			"USER_NAME":       sshUser,
			"LOG_STDOUT":      "true",
		},

		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort("2222/tcp").WithStartupTimeout(10 * time.Second),
		Name:         "test-server",
		HostConfigModifier: func(cfg *container.HostConfig) {
			cfg.PortBindings = network.PortMap{
				network.MustParsePort("2222"): []network.PortBinding{{
					HostIP:   netip.MustParseAddr("0.0.0.0"),
					HostPort: "2222",
				}},
			}
		},
		LifecycleHooks: []testcontainers.ContainerLifecycleHooks{{
			PostCreates: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PostCreates")
					return nil
				},
			},
			PreStarts: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PreStarts")
					return nil
				},
			},
			PostStarts: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PostStarts")
					return nil
				},
			},
			PostReadies: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PostReadies")
					return nil
				},
			},
			PreStops: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PreStops")
					return nil
				},
			},
			PostStops: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PostStops")
					return nil
				},
			},
			PreTerminates: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PreTerminates")
					return nil
				},
			},
			PostTerminates: []testcontainers.ContainerHook{
				func(ctx context.Context, ctr testcontainers.Container) error {
					t.Log(" -> PostTerminates")
					return nil
				},
			},
		}},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      nil,
			Consumers: []testcontainers.LogConsumer{log},
		},
	}
}
