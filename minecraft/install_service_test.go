package minecraft

import (
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstaller_DownloadServer(t *testing.T) {
	t.Run("", func(t *testing.T) {
		c := NewInstallService(WithTimeout(1 * time.Second))

		v := versions.VersionInfoResponse{
			Downloads: versions.Downloads{
				Server: versions.Artifact{
					URL:  "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar",
					SHA1: "15c777e2cfe0556eef19aab534b186c0c6f277e1",
				},
			},
		}
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		serverFile, err := c.DownloadServer(v, dest)
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(dest, "server.jar"), serverFile)
	})
}
