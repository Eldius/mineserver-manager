package minecraft

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapBackupFiles(t *testing.T) {
	t.Run("given a folder containing some valid backup files should return the right size files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./test_samples/")
		assert.NoError(t, err)
		assert.Len(t, files, 31)
	})
	t.Run("given a folder without valid backup files should return an empty files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./config/")
		assert.NoError(t, err)
		assert.Empty(t, files)
	})
}
