package config

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
