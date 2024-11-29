package utils

import (
	"context"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFileName(t *testing.T) {
	t.Run("given java runtime version 17 for linux in a platform amd64 should return the right file name", func(t *testing.T) {
		name := GetFileName("https://aka.ms/download-jdk/microsoft-jdk-17.0.7-linux-x64.tar.gz")

		assert.Equal(t, "microsoft-jdk-17.0.7-linux-x64.tar.gz", name)
	})

	t.Run("given java runtime version 17 for linux in a platform arm64 should return the right file name", func(t *testing.T) {
		name := GetFileName("https://aka.ms/download-jdk/microsoft-jdk-17.0.7-linux-aarch64.tar.gz")

		assert.Equal(t, "microsoft-jdk-17.0.7-linux-aarch64.tar.gz", name)
	})

	t.Run("given java runtime version 11 for linux in a platform arm64 should return the right file name", func(t *testing.T) {
		name := GetFileName("https://aka.ms/download-jdk/microsoft-jdk-11.0.19-linux-aarch64.tar.gz")

		assert.Equal(t, "microsoft-jdk-11.0.19-linux-aarch64.tar.gz", name)
	})
}

func TestDownloadFile(t *testing.T) {
	gock.New("https://piston-data.mojang.com").
		Get("/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar").
		Reply(200).
		File("../minecraft/mojang/samples/server.zip")

	outDir, err := os.MkdirTemp(os.TempDir(), "mineserver-testing-*")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(outDir)
	}()

	destFile := filepath.Join(outDir, "server.zip")

	assert.NoError(t, DownloadFile(context.Background(), 200*time.Millisecond, "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar", destFile))
}

func TestValidateFileIntegrity(t *testing.T) {
	t.Run("given a valid file hash should return success", func(t *testing.T) {
		file := "../minecraft/mojang/samples/server.zip"
		assert.NoError(t, ValidateFileIntegrity(context.Background(), file, "fe5c3e7c6983ac7ea8a23bd9f2d8b235128633e8"))
	})

	t.Run("given a valid file hash should return error", func(t *testing.T) {
		file := "../minecraft/mojang/samples/server.zip"
		assert.Error(t, ValidateFileIntegrity(context.Background(), file, "fe5c3e7c6983ac7ea8a23bd9f2d8b2351286abcd"))
	})
}
