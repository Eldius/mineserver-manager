package provisioner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateScript(t *testing.T) {
	t.Run("headless default", func(t *testing.T) {
		script, err := StartScript(
			WithJDKPath("/tmp/java/jdk"),
			WithMemLimit("1G"),
			WithServerFile("my-server-123.jar"),
			WithHeadless(true),
		)
		assert.Nil(t, err)
		assert.Contains(t, script, "#!/bin/bash")
		assert.Contains(t, script, "/tmp/java/jdk/java")
		assert.Contains(t, script, "-Xms1G")
		assert.Contains(t, script, "-Xmx1G")
		assert.Contains(t, script, "-jar my-server-123.jar")
		assert.Contains(t, script, "--nogui")
	})
	t.Run("not headless default", func(t *testing.T) {
		script, err := StartScript(
			WithJDKPath("/tmp/java/jdk-123"),
			WithMemLimit("500m"),
			WithServerFile("my-server-456.jar"),
			WithHeadless(false),
		)
		assert.Nil(t, err)
		assert.Contains(t, script, "#!/bin/bash")
		assert.Contains(t, script, "/tmp/java/jdk-123/java")
		assert.Contains(t, script, "-Xms500m")
		assert.Contains(t, script, "-Xmx500m")
		assert.Contains(t, script, "-jar my-server-456.jar")
		assert.NotContains(t, script, "--nogui")
	})
}

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
