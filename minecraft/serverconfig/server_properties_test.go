package serverconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDefaultServerProperties(t *testing.T) {

	p, err := GetDefaultServerProperties()
	assert.Nil(t, err)

	assert.True(t, p.AllowNether)
	assert.False(t, p.AllowFlight)
}
