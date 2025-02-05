package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScriptParams_ToScript(t *testing.T) {
	t.Run("given an script configuration without log configuration file should create a start script without log configuration", func(t *testing.T) {
		s, err := StartScript(
			WithMemLimit("512m"))
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)
		assert.Contains(t, s, "-Xms512m")
		assert.Contains(t, s, "-Xmx512m")
		assert.Contains(t, s, "-jar server.jar")
		assert.NotContains(t, s, "-Dlog4j.configurationFile")
		assert.NotContains(t, s, "--nogui")
	})

	t.Run("given an script configuration with log configuration file should create a start script with log configuration", func(t *testing.T) {
		s, err := StartScript(
			WithMemLimit("4g"),
			WithLogConfigFile(true),
		)
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)

		assert.Contains(t, s, "-Xms4g")
		assert.Contains(t, s, "-Xmx4g")
		assert.Contains(t, s, "-jar server.jar")
		assert.Contains(t, s, "-Dlog4j.configurationFile=${INSTALL_PATH}/log4j2.xml")
		assert.NotContains(t, s, "--nogui")
	})

	t.Run("given an script configuration without log configuration and headless mode enabled file should create a start script without log configuration and nogui parameter", func(t *testing.T) {
		s, err := StartScript(
			WithMemLimit("512m"),
			WithHeadless(true),
		)
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)

		assert.Contains(t, s, "-Xms512m")
		assert.Contains(t, s, "-Xmx512m")
		assert.Contains(t, s, "-jar server.jar")
		assert.NotContains(t, s, "-Dlog4j.configurationFile")
		assert.Contains(t, s, "--nogui")
	})
}

func TestWithServerFlavour(t *testing.T) {
	t.Run("given an empty flavour should return an empty attribute", func(t *testing.T) {
		f := InstanceOpts{}
		WithServerFlavour("")(&f)

		assert.Equal(t, EmptyServerSoftware, f.Flavor)
	})

	t.Run("given a valid flavour should return an equal value from attribute", func(t *testing.T) {
		f := InstanceOpts{}
		WithServerFlavour(string(PurpurServerSoftware))(&f)

		assert.Equal(t, PurpurServerSoftware, f.Flavor)

		f1 := InstanceOpts{}
		WithServerFlavour(string(EmptyServerSoftware))(&f1)

		assert.Equal(t, EmptyServerSoftware, f1.Flavor)

		f2 := InstanceOpts{}
		WithServerFlavour(string(VanillaServerSoftware))(&f2)

		assert.Equal(t, VanillaServerSoftware, f2.Flavor)
	})

	t.Run("given an invalid flavour should return an empty attribute", func(t *testing.T) {
		f := InstanceOpts{}
		WithServerFlavour("some random value")(&f)

		assert.Equal(t, EmptyServerSoftware, f.Flavor)
	})
}
