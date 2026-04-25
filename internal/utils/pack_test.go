package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackedFileName(t *testing.T) {
	assert.Equal(t,
		"mine/server.properties",
		packedFileName("/home/my-user/mine/server.properties", "/home/my-user/"),
	)
	assert.Equal(t,
		"mine/server.properties",
		packedFileName("/home/my-user/mine/server.properties", "/home/my-user"),
	)
	assert.Equal(t,
		"mine/sub-folder/server.properties",
		packedFileName("/home/my-user/mine/sub-folder/server.properties", "/home/my-user"),
	)
}
