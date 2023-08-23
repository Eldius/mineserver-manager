package versions

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestListVersions(t *testing.T) {
	t.Run("given a simple version query should return success with the right data", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://launchermeta.mojang.com").
			Get("/mc/game/version_manifest.json").
			Reply(200).
			File("./samples/versions.json")

		gock.New("https://piston-meta.mojang.com").
			Get("/v1/packages/7efb232e2903bea16d7bf7b4a5ea768453cf92ea/1.20.json").
			Reply(200).
			File("./samples/1.20.json")

		c := NewClient(WithTimeout(1 * time.Second))

		v, err := c.ListVersions()
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, 696, len(v.Versions))

		assert.Equal(t, "1.20", v.Latest.Release)
		assert.Equal(t, "1.20.1-rc1", v.Latest.Snapshot)

		r, err := v.GetLatestRelease()
		assert.Nil(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "1.20", r.ID)
	})

	t.Run("given an specific version should return its info", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://launchermeta.mojang.com").
			Get("/mc/game/version_manifest.json").
			Reply(200).
			File("./samples/versions.json")

		gock.New("https://piston-meta.mojang.com").
			Get("/v1/packages/7efb232e2903bea16d7bf7b4a5ea768453cf92ea/1.20.json").
			Reply(200).
			File("./samples/1.20.json")

		c := NewClient(WithTimeout(1 * time.Second))

		v, err := c.ListVersions()
		assert.Nil(t, err)
		assert.NotNil(t, v)

		lr, err := v.GetLatestRelease()
		assert.Nil(t, err)
		assert.NotNil(t, lr)

		info, err := c.GetVersionInfo(*lr)
		assert.Nil(t, err)
		assert.NotNil(t, info)

		assert.Equal(t, "https://piston-data.mojang.com/v1/objects/e575a48efda46cf88111ba05b624ef90c520eef1/client.jar", info.Downloads.Client.URL)
		assert.Equal(t, "https://piston-data.mojang.com/v1/objects/15c777e2cfe0556eef19aab534b186c0c6f277e1/server.jar", info.Downloads.Server.URL)

		assert.Equal(t, "1.20", info.ID)
		assert.Equal(t, 17, info.JavaVersion.MajorVersion)
	})

	t.Run("given an specific version should return its info", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://launchermeta.mojang.com").
			Get("/mc/game/version_manifest.json").
			Reply(200).
			Delay(5 * time.Second).
			File("./samples/versions.json")

		c := NewClient(WithTimeout(1 * time.Second))

		v, err := c.ListVersions()
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})
}
