package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
