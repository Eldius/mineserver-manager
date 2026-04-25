package installer

import (
	"context"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloader_DownloadServer(t *testing.T) {
	defer gock.Off()
	t.Run("given a version with the right checksum value Name should download and validate file checksum with success", func(t *testing.T) {
		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("../mojang/samples/server.zip")

		ctx := context.Background()

		d := NewDownloader(1 * time.Second)

		url := "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar"
		sha1 := "fe5c3e7c6983ac7ea8a23bd9f2d8b235128633e8"
		
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		defer os.RemoveAll(dest)

		serverFile, err := d.DownloadServer(ctx, url, sha1, dest)
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(dest, "server.jar"), serverFile)

		if stat, err := os.Stat(serverFile); !assert.Nil(t, err) || !assert.False(t, stat.IsDir()) {
			assert.FailNow(t, "invalid server.jar file")
		}
	})

	t.Run("given a version with the wrong checksum value Name should download and validate file checksum without success", func(t *testing.T) {
		gock.New("https://piston-data.mojang.com").
			Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
			Reply(200).
			File("../mojang/samples/server.zip")

		ctx := context.Background()

		d := NewDownloader(1 * time.Second)

		url := "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar"
		sha1 := "2e49d5731f612a27506fc777ee146fc4080312de"
		
		dest, err := os.MkdirTemp(os.TempDir(), "mine-test-*")
		assert.Nil(t, err)
		defer os.RemoveAll(dest)

		serverFile, err := d.DownloadServer(ctx, url, sha1, dest)
		assert.NotNil(t, err)
		assert.Empty(t, serverFile)
	})
}
