package java

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDownloadErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("given an unsupported version should return ErrVersionNotSupported", func(t *testing.T) {
		_, err := Download(ctx, 99, "amd64", "linux", time.Second)
		assert.ErrorIs(t, err, ErrVersionNotSupported)
	})

	t.Run("given an unsupported platform should return ErrPlatformNotSupported", func(t *testing.T) {
		_, err := Download(ctx, 17, "amd64", "windows", time.Second)
		assert.ErrorIs(t, err, ErrPlatformNotSupported)
	})

	t.Run("given an unsupported arch should return ErrPlatformNotSupported", func(t *testing.T) {
		_, err := Download(ctx, 17, "ppc64", "linux", time.Second)
		assert.ErrorIs(t, err, ErrPlatformNotSupported)
	})
}
