package minecraft

import (
	"github.com/eldius/mineserver-manager/minecraft/versions"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstaller_DownloadServer(t *testing.T) {
	t.Run("", func(t *testing.T) {
		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("./versions/samples/versions.json")

		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("./versions/samples/server.zip")

		c := NewInstallService(WithTimeout(1 * time.Second))

		v := versions.VersionInfoResponse{
			Downloads: versions.Downloads{
				Server: versions.Artifact{
					URL:  "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar",
					SHA1: "2e49d5731f612a27506fc777ee146fc4080312de",
				},
			},
		}
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		serverFile, err := c.DownloadServer(v, dest)
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(dest, "server.jar"), serverFile)

		if stat, err := os.Stat(serverFile); assert.Nil(t, err) {
			assert.False(t, stat.IsDir())
		}

	})
}