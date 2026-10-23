package repository

import (
	"context"
	"github.com/eldius/mineserver-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestStormRepository(t *testing.T) {
	repo, err := NewStormRepository(filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	defer func() {
		_ = repo.Close()
	}()

	ctx := context.Background()

	t.Run("getting a missing instance returns ErrNotFound", func(t *testing.T) {
		_, err := repo.GetInstance(ctx, "missing-id")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("deleting a missing instance returns ErrNotFound", func(t *testing.T) {
		err := repo.DeleteInstance(ctx, "missing-id")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("saved instance can be retrieved", func(t *testing.T) {
		inst := model.NewInstance("test-instance", "/tmp/test-instance", model.ServerProperties{})
		require.NoError(t, repo.SaveInstance(ctx, inst))

		got, err := repo.GetInstance(ctx, inst.ID)
		require.NoError(t, err)
		assert.Equal(t, inst.ID, got.ID)
		assert.Equal(t, "test-instance", got.Name)
	})
}
