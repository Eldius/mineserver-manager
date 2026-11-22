package minecraft

import (
	"context"
	"testing"
	"time"

	"github.com/eldius/cli-telemetry-helper/pkg/setup"
	"github.com/stretchr/testify/assert"
)

func initTest(t *testing.T) {
	t.Helper()
	assert.Nil(t, setup.InitViper(t.Context(), setup.Service{
		Name:    "my-testing-case",
		Version: "",
	}, setup.WithConfigFileLookupDirectories("../..")))
}

func TestRestoreBackupFileNotFound(t *testing.T) {
	initTest(t)

	s := NewBackupService()
	err := s.Restore(context.Background(), t.TempDir(), "/nonexistent/my-backup.zip")
	assert.ErrorIs(t, err, ErrBackupFileNotFound)
}

func TestMapBackupFiles(t *testing.T) {
	initTest(t)

	t.Run("given a folder containing some valid backup files should return the right size files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./testdata/")
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Len(t, files["mybackup_file"], 5)
		assert.Len(t, files["my_other_backup_file"], 1)

		myBkpFileOlder := files["mybackup_file"].olderFile()
		ts, _ := time.Parse(bkpTimestampFormat, "2024-12-29_00-00-01")

		assert.NotNil(t, myBkpFileOlder)

		assert.Equal(t, ts, myBkpFileOlder.Timestamp)
		assert.Equal(t, "testdata/mybackup_file_2024-12-29_00-00-01_backup.zip", myBkpFileOlder.Path)
		assert.Equal(t, "mybackup_file", myBkpFileOlder.Name)
	})
	t.Run("given a folder without valid backup files should return an empty files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./config/")
		assert.NoError(t, err)
		assert.Empty(t, files)
	})
}
