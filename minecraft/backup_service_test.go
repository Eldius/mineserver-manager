package minecraft

import (
	"context"
	"github.com/eldius/initial-config-go/configs"
	"github.com/eldius/initial-config-go/setup"
	"github.com/eldius/mineserver-manager/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	_ = setup.InitSetup(config.AppName, setup.WithDefaultValues(map[string]any{
		configs.LogLevelKey:  configs.LogLevelDEBUG,
		configs.LogFormatKey: configs.LogFormatJSON,
	}))
}

func TestMapBackupFiles(t *testing.T) {
	t.Run("given a folder containing some valid backup files should return the right size files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./test_samples/")
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Len(t, files["mybackup_file"], 5)
		assert.Len(t, files["my_other_backup_file"], 1)

		myBkpFileOlder := files["mybackup_file"].olderFile()
		ts, _ := time.Parse(bkpTimestampFormat, "2024-12-29_00-00-01")

		assert.NotNil(t, myBkpFileOlder)

		assert.Equal(t, ts, myBkpFileOlder.Timestamp)
		assert.Equal(t, "test_samples/mybackup_file_2024-12-29_00-00-01_backup.zip", myBkpFileOlder.Path)
		assert.Equal(t, "mybackup_file", myBkpFileOlder.Name)
	})
	t.Run("given a folder without valid backup files should return an empty files map", func(t *testing.T) {
		files, err := mapBackupFiles(context.Background(), "./config/")
		assert.NoError(t, err)
		assert.Empty(t, files)
	})
}
