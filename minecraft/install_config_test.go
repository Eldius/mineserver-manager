package minecraft

import (
	"github.com/eldius/mineserver-manager/minecraft/serverconfig"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScriptParams_ToScript(t *testing.T) {
	t.Run("given an script configuration without log configuration file should create a start script without log configuration", func(t *testing.T) {
		sp := &serverconfig.StartupParams{
			Xmx: "512m",
			Xms: "256m",
		}
		s, err := sp.ToScript()
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)
		assert.Contains(t, s, "-Xms256m")
		assert.Contains(t, s, "-Xmx512m")
		assert.Contains(t, s, "-jar server.jar")
		assert.NotContains(t, s, "-Dlog4j.configurationFile")
		assert.NotContains(t, s, "--nogui")
	})

	t.Run("given an script configuration with log configuration file should create a start script with log configuration", func(t *testing.T) {
		sp := &serverconfig.StartupParams{
			Xmx:           "4g",
			Xms:           "4g",
			LogConfigFile: "log4j2.xml",
		}

		s, err := sp.ToScript()
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)

		assert.Contains(t, s, "-Xms4g")
		assert.Contains(t, s, "-Xmx4g")
		assert.Contains(t, s, "-jar server.jar")
		assert.Contains(t, s, "-Dlog4j.configurationFile=log4j2.xml")
		assert.NotContains(t, s, "--nogui")
	})

	t.Run("given an script configuration without log configuration and headless mode enabled file should create a start script without log configuration and nogui parameter", func(t *testing.T) {
		sp := &serverconfig.StartupParams{
			Xmx:      "512m",
			Xms:      "256m",
			Headless: true,
		}
		s, err := sp.ToScript()
		t.Logf("s: '%v'", s)
		assert.Nil(t, err)

		assert.Contains(t, s, "-Xms256m")
		assert.Contains(t, s, "-Xmx512m")
		assert.Contains(t, s, "-jar server.jar")
		assert.NotContains(t, s, "-Dlog4j.configurationFile")
		assert.Contains(t, s, "--nogui")
	})
}