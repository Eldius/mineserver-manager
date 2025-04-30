package utils

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"path/filepath"
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

	ctx, cancel1 := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel1()

	//assert.NoError(t, Execute(ctx, "localhost", "2222"))
	assert.NoError(t, Execute(ctx, ip, "2222"))

	//time.Sleep(30 * time.Second)
}

//type logConsumer struct {
//	t *testing.T
//}
//
//func (c *logConsumer) Accept(log testcontainers.Log) {
//	c.t.Log(log)
//}

type myLogConsumer struct {
	// You can add fields here to store or process the logs
	testcontainers.LogConsumer
	t *testing.T
}

type myWriter struct {
	io.Writer
	t *testing.T
}

func (c *myLogConsumer) Accept(log testcontainers.Log) {
	// Handle the log line here
	c.t.Logf("[container] %s", log.Content)
}

func (w *myWriter) Write(p []byte) (n int, err error) {
	w.t.Logf("[build] %s", p)
	return len(p), nil
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

func containerRequestFromDockerfile(t *testing.T) testcontainers.ContainerRequest {
	t.Helper()

	w := &myWriter{
		t: t,
	}
	log := &myLogConsumer{
		t: t,
	}
	return testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:              filepath.Join(".", "testing", "ssh"),
			Dockerfile:           "Dockerfile",
			BuildLogWriter:       w,
			KeepImage:            false,
			BuildOptionsModifier: nil,
		},
		Env: map[string]string{
			"PUID": "1000",
			"PGID": "1000",
			"TZ":   "Etc/UTC",
			//"PUBLIC_KEY":         "yourpublickey",
			//"PUBLIC_KEY_FILE":    "/path/to/file",
			//"PUBLIC_KEY_DIR":     "/path/to/directory/containing/_only_/pubkeys",
			//"PUBLIC_KEY_URL":  "https://github.com/username.keys",
			"SUDO_ACCESS":     "false",
			"PASSWORD_ACCESS": "true",
			"USER_PASSWORD":   sshPass,
			//"USER_PASSWORD_FILE": "/path/to/file",
			"USER_NAME":  sshUser,
			"LOG_STDOUT": "true",
		},
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(exposedPort),
		Name:         "test-server",
		HostConfigModifier: func(cfg *container.HostConfig) {
			//cfg.
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
			//"USER_PASSWORD_FILE": "/path/to/file",
			//"PUBLIC_KEY":         "yourpublickey",
			//"PUBLIC_KEY_FILE":    "/path/to/file",
			//"PUBLIC_KEY_DIR":     "/path/to/directory/containing/_only_/pubkeys",
			//"PUBLIC_KEY_URL":  "https://github.com/username.keys",
		},

		ExposedPorts: []string{exposedPort},
		//WaitingFor:   wait.ForListeningPort(exposedPort),
		WaitingFor: wait.ForLog("Server listening on 0.0.0.0 port 2222"),
		Name:       "test-server",
		HostConfigModifier: func(cfg *container.HostConfig) {
			cfg.PortBindings = nat.PortMap{
				"2222": []nat.PortBinding{{
					HostIP:   "0.0.0.0",
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
