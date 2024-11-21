package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
