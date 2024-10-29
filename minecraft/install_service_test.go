package minecraft

import (
	"context"
	"github.com/eldius/mineserver-manager/minecraft/mojang"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstaller_DownloadServer(t *testing.T) {
	t.Run("given a version with the right checksum value name should download and validate file checksum with success", func(t *testing.T) {
		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("./mojang/samples/server.zip")

		ctx := context.Background()

		c := NewInstallService(WithTimeout(1 * time.Second))

		v := mojang.VersionInfoResponse{
			Downloads: mojang.Downloads{
				Server: mojang.Artifact{
					URL:  "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar",
					SHA1: "fe5c3e7c6983ac7ea8a23bd9f2d8b235128633e8",
				},
			},
		}
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		serverFile, err := c.DownloadServer(ctx, v, dest)
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(dest, "server.jar"), serverFile)

		if stat, err := os.Stat(serverFile); !assert.Nil(t, err) || !assert.False(t, stat.IsDir()) {
			assert.FailNow(t, "invalid server.jar file")
		}
	})

	t.Run("given a version with the wrong checksum value name should download and validate file checksum without success", func(t *testing.T) {
		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("./mojang/samples/server.zip")

		ctx := context.Background()

		c := NewInstallService(WithTimeout(1 * time.Second))

		v := mojang.VersionInfoResponse{
			Downloads: mojang.Downloads{
				Server: mojang.Artifact{
					URL:  "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar",
					SHA1: "2e49d5731f612a27506fc777ee146fc4080312de",
				},
			},
		}
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		serverFile, err := c.DownloadServer(ctx, v, dest)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrChecksumValidationFailed)
		assert.Empty(t, serverFile)
	})
}

func TestInstaller_CreateStartScript(t *testing.T) {
	t.Run("", func(t *testing.T) {

	})
}
